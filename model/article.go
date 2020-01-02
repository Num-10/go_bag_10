package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Article struct {
	Content       string `gorm:"column:content"`
	CoverImageURL string `gorm:"column:cover_image_url"`
	Created       int    `gorm:"column:created"`
	Desc          string `gorm:"column:desc"`
	ID            int    `gorm:"column:id;primary_key"`
	IsMarrow      int    `gorm:"column:is_marrow"`
	IsTop         int    `gorm:"column:is_top"`
	Sort          int    `gorm:"column:sort"`
	Status        int    `gorm:"column:status"`
	TagID         int    `gorm:"column:tag_id"`
	Title         string `gorm:"column:title"`
	Updated       int    `gorm:"column:updated"`
	UserID        int    `gorm:"column:user_id"`
}

// TableName sets the insert table name for this struct type
func (a *Article) TableName() string {
	return "blog_article"
}

func (a *Article) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("Created", time.Now().Unix())
	return nil
}

func (a *Article) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("Updated", time.Now().Unix())
	return nil
}

