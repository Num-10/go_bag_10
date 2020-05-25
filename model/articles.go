package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Articles struct {
	ID             int    `gorm:"column:id;primary_key"`
	ColumnID       int    `gorm:"column:column_id"`
	ArticleID      int    `gorm:"column:article_id"`
	ArticleTitle   string `gorm:"column:article_title"`
	ArticleSummary string `gorm:"column:article_summary"`
	UpdatedTime    int    `gorm:"column:updated_time"`
}

// TableName sets the insert table name for this struct type
func (a *Articles) TableName() string {
	return "articles"
}

func (a *Articles) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedTime", time.Now().Unix())
	return nil
}

func (a *Articles) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedTime", time.Now().Unix())
	return nil
}

func (a *Articles) Create() error {
	return Db.Create(a).Error
}

func (a *Articles) Find(where interface{}, order string) {
	Db.Where(where).Order(order).First(a)
}

func (a *Articles) Update(where, data interface{}) error {
	return Db.Model(&Articles{}).Where(where).Update(data).Error
}

func (a *Articles) Delete(where interface{}) error {
	return Db.Where(where).Delete(Articles{}).Error
}

func (a *Articles) GetList(where map[string]interface{}, extra map[string]interface{}, articles interface{}, count *int) {
	query := Db.Model(&Articles{}).Where(where)
	if _, ok := extra["field"]; ok {
		query = query.Select(extra["field"])
	}
	if _, ok := extra["multi_like_search"]; ok && extra["multi_like_search"] != "" {
		extra["multi_like_search"] = "%" + (extra["multi_like_search"]).(string) + "%"
		query = query.Where("`title` like ? or `desc` like ? or `content` like ?", extra["multi_like_search"], extra["multi_like_search"], extra["multi_like_search"])
	}
	if _, ok := extra["group"]; ok {
		query = query.Group((extra["group"]).(string))
	}
	if _, ok := extra["count"]; ok {
		query = query.Count(count)
		return
	}
	if _, ok := extra["order"]; ok {
		query = query.Order(extra["order"])
	}
	page, ok := extra["page"]
	pageSize, pok := extra["page_size"]
	if ok && pok {
		query = query.Limit(pageSize).Offset(((page).(int) - 1) * (pageSize).(int))
	}
	query = query.Scan(articles)
	return
}
