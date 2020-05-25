package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Columns struct {
	ID                int     `gorm:"column:id;primary_key"`
	ColumnID          int     `gorm:"column:column_id"`
	ColumnTitle       string  `gorm:"column:column_title"`
	ColumnSubtitle    string  `gorm:"column:column_subtitle"`
	ColumnType        int     `gorm:"column:column_type"`
	ColumnPrice       float64 `gorm:"column:column_price"`
	ColumnPriceMarket float64 `gorm:"column:column_price_market"`
	ColumnBeginTime   int     `gorm:"column:column_begin_time"`
	ColumnEndTime     int     `gorm:"column:column_end_time"`
	ColumnSku         int     `gorm:"column:column_sku"`
	ColumnCoverInner  string  `gorm:"column:column_cover_inner"`
	ColumnCoverWxlite string  `gorm:"column:column_cover_wxlite"`
	AuthorName        string  `gorm:"column:author_name"`
	AuthorIntro       string  `gorm:"column:author_intro"`
	ArticleDoneCount  int     `gorm:"column:article_done_count"`
	ArticleTotalCount int     `gorm:"column:article_total_count"`
	UpdatedTime       int     `gorm:"column:updated_time"`
}

// TableName sets the insert table name for this struct type
func (a *Columns) TableName() string {
	return "columns"
}

func (a *Columns) BeforeCreate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedTime", time.Now().Unix())
	return nil
}

func (a *Columns) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedTime", time.Now().Unix())
	return nil
}

func (a *Columns) ColumnAdd() error {
	return Db.Create(a).Error
}

func (a *Columns) Find(where interface{}, order string) {
	Db.Where(where).Order(order).First(a)
}

func (a *Columns) Update(where, data interface{}) error {
	return Db.Model(&Columns{}).Where(where).Update(data).Error
}

func (a *Columns) GetList(where map[string]interface{}, extra map[string]interface{}, articles interface{}, count *int) {
	query := Db.Model(&Columns{}).Where(where)
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
