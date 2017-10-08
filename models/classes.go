package models

import (
	"github.com/jinzhu/gorm"
)

type Class struct {
	gorm.Model
	Name        string
	Description string
	Videos      []Video
}

type ClassDB interface {
	GetAll() ([]Class, error)
	GetClass(name string) (*Class, error)
	CreateClass(class *Class) error
}

type ClassService interface {
	ClassDB
}

func NewClassService(db *gorm.DB) ClassService {
	return &classService{
		ClassDB: &classGorm{db},
	}
}

type classService struct {
	ClassDB
}

type classGorm struct {
	db *gorm.DB
}

func (cg *classGorm) CreateClass(class *Class) error {
	return cg.db.Create(class).Error
}

func (cg *classGorm) GetAll() ([]Class, error) {
	classes := []Class{}
	if err := cg.db.Find(&classes).Error; err != nil {
		return nil, err
	}
	return classes, nil
}

func (cg *classGorm) GetClass(name string) (*Class, error) {
	class := Class{}
	err := cg.db.Where("name = ?", name).First(&class).Error
	if err != nil {
		return nil, err
	}
	return &class, nil
}
