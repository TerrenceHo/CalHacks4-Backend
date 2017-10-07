package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Services struct {
	db   *gorm.DB
	User UserService
}

type ServicesConfig func(*Services) error

func WithGorm(dialect, connectionInfo string) ServicesConfig {
	return func(s *Services) error {
		db, err := gorm.Open(dialect, connectionInfo)
		if err != nil {
			return err
		}
		s.db = db
		return nil
	}
}

func WithLogMode(mode bool) ServicesConfig {
	return func(s *Services) error {
		s.db.LogMode(mode)
		return nil
	}
}

func WithUser(pepper string) ServicesConfig {
	return func(s *Services) error {
		s.User = NewUserService(s.db, pepper)
		return nil
	}
}

func NewServices(cfgs ...ServicesConfig) (*Services, error) {
	var s Services
	for _, cfg := range cfgs {
		if err := cfg(&s); err != nil {
			return nil, err
		}
	}
	return &s, nil
}

// Closes connection to database
func (s *Services) Close() error {
	return s.db.Close()
}

func (s *Services) DestructiveReset() error {
	err := s.db.DropTableIfExists(&User{}).Error
	if err != nil {
		return nil
	}
	return s.AutoMigrate()
}

// Attempts to migrate User, InboundVehicle, and OutboundVehicle
func (s *Services) AutoMigrate() error {
	return s.db.AutoMigrate(&User{}).Error
}
