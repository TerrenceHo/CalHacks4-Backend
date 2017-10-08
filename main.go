package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/TerrenceHo/CalHacks4-Backend/config"
	"github.com/TerrenceHo/CalHacks4-Backend/controllers"
	"github.com/TerrenceHo/CalHacks4-Backend/models"
	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()
	// connection := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", host, port, user, name)
	services, err := models.NewServices(
		models.WithGorm(cfg.DatabaseDialect(), cfg.DatabaseConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper),
		models.WithClass(),
		models.WithVideo(),
	)
	must(err)
	defer services.Close()
	err = services.AutoMigrate()
	must(err)

	usersC := controllers.NewUsers(services.User, cfg.SignKey)
	classesC := controllers.NewClasses(services.Class, services.Video)

	router := mux.NewRouter()
	router.HandleFunc("/", homePage).Methods("GET")
	router.HandleFunc("/api/v1/user/register", usersC.Create).Methods("POST")
	router.HandleFunc("/api/v1/user/login", usersC.Login).Methods("POST")

	router.HandleFunc("/api/v1/classes", classesC.GetAllClasses).Methods("GET")
	router.HandleFunc("/api/v1/classes/{id}", classesC.GetClass).Methods("GET")
	router.HandleFunc("/api/v1/classes/create", classesC.Create).Methods("POST")
	router.HandleFunc("/api/v1/classes/upload", classesC.Upload).Methods("POST")
	router.HandleFunc("/api/v1/classes/search", classesC.GetByKeyword).Methods("POST")

	log.Println("Listening on Port", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<h1>Hello World!</h1>")
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
