package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Video struct {
	gorm.Model
	ClassID           uint
	URL               string
	Topics            pq.StringArray `gorm:"type:varchar(200)[]"`
	Related_Resources pq.StringArray `gorm:"type:varchar(200)[]"`
}

type VideoDB interface {
	GetAll(id uint) ([]Video, error)
	Create(video *Video) error
	GetByKeyword(id uint, keyword string) ([]Video, error)
}

type VideoService interface {
	VideoDB
}

func NewVideoService(db *gorm.DB) VideoService {
	return &videoService{
		VideoDB: &videoGorm{db},
	}
}

type videoService struct {
	VideoDB
}

type videoGorm struct {
	db *gorm.DB
}

func (vg *videoGorm) Create(video *Video) error {
	return vg.db.Create(video).Error
}

func (vg *videoGorm) GetAll(id uint) ([]Video, error) {
	videos := []Video{}
	if err := vg.db.Where("class_id = ?", id).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

func (vg *videoGorm) GetByKeyword(id uint, keyword string) ([]Video, error) {
	videos := []Video{}
	db := vg.db.Where("class_id = ?", id)
	err := db.Where("? = ANY(topics)", keyword).Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}
