package main

import (
	"encoding/json"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	dbPath   = "data/db.sqlite3"
	rectsKey = "rects"
)

type Menu struct {
	PublishDate  time.Time `gorm:"primary_key"`
	DeliveryDate time.Time
	Items        string
}

type KeyValue struct {
	Key   string `gorm:"primarykey"`
	Value string
}

var Db *gorm.DB

func InitDb() {
	if db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{}); err != nil {
		log.Fatal(err)
		return
	} else {
		Db = db
	}

	if err := Db.AutoMigrate(&Menu{}); err != nil {
		log.Fatal(err)
		return
	}
	if err := Db.AutoMigrate(&KeyValue{}); err != nil {
		log.Fatal(err)
		return
	}
}

func GetRects() ([]FloatRect, error) {
	var kv KeyValue
	Db.Find(&kv, "key = ?", rectsKey)
	var rects []FloatRect
	if err := json.Unmarshal([]byte(kv.Value), &rects); err != nil {
		return nil, err
	}
	return rects, nil
}

func SaveRects(rects string) {
	Db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(&KeyValue{
		Key:   rectsKey,
		Value: rects,
	})
}

func SaveMenu(menu Menu) {
	Db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "publish_date"}},
		DoUpdates: clause.AssignmentColumns([]string{"items", "delivery_date"}),
	}).Create(&menu)
}

func GetMenu(date time.Time) (Menu, error) {
	var menu Menu
	if err := Db.Where("publish_date = ?", date).First(&menu).Error; err != nil {
		return menu, err
	}
	return menu, nil
}
