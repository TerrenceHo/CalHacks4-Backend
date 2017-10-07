package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type PostgresConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (c *Config) DatabaseDialect() string {
	return "postgres"
}

func (c *Config) DatabaseConnectionInfo() string {
	if c.Env == "prod" {
		return os.Getenv("DATABASE_URL")
	} else {
		if c.Database.Password == "" {
			return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
				c.Database.Host, c.Database.Port, c.Database.User, c.Database.Name)
		}
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			c.Database.Host, c.Database.Port, c.Database.User, c.Database.Password, c.Database.Name)

	}
}

type Config struct {
	Port     string         `json:"port"`
	Env      string         `json:"env"`
	Pepper   string         `json:"pepper"`
	Database PostgresConfig `json:"database"`

	PubKeyPath            string `json:"pubKeyPath"`
	PrivKeyPath           string `json:"privKeyPath"`
	PassResetSecretString string `json:"passResetSecret"`

	VerifyKey []byte
	SignKey   []byte
}

func LoadConfig() *Config {
	c := readJSONConfig()
	c.checkProd()
	c.loadJWTKeys()

	fmt.Println("Successfully Loaded Config File")
	return c
}

func (c *Config) IsProd() bool {
	return c.Env == "prod"
}

func readJSONConfig() *Config {
	f, err := os.Open("config/.config")
	if err != nil {
		log.Fatal("No config file", err)
	}
	var c Config
	dec := json.NewDecoder(f)
	err = dec.Decode(&c)
	if err != nil {
		log.Fatal("JSON malformed", err)
	}
	return &c
}

func (c *Config) checkProd() {
	if os.Getenv("SERVER_ENV") == "production" || os.Getenv("SERVER_ENV") == "staging" {
		c.Port = os.Getenv("PORT")
		c.Env = "prod"
	}
}

func (c *Config) loadJWTKeys() {
	// Load JWT keys
	var err error
	c.SignKey, err = ioutil.ReadFile(c.PrivKeyPath)
	if err != nil {
		log.Fatal("Error reading private key:", err)
	}
	c.VerifyKey, err = ioutil.ReadFile(c.PubKeyPath)
	if err != nil {
		log.Fatal("Error reading public key:", err)
	}
}

// func (c *Config) loadPassReset() {
// 	// Load PasswordReset Secret
// 	c.PassResetSecret = []byte(c.PassResetSecretString)
// }
