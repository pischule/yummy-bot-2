package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
	"yummy-bot/ocr"
)

const rectsFilePath = "data/rects.json"
const menuFilePath = "data/menu.json"

type Menu struct {
	PublishDate  time.Time `json:"publish_date"`
	DeliveryDate time.Time `json:"delivery_date"`
	Items        []string  `json:"items"`
}

func readObject[T any](v T, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func writeObject[T any](v T, path string) error {
	f, err := os.Create(path)
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Println("error while closing file: ", err)
		}
	}(f)
	if err != nil {
		return err
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

func ReadRects() ([]ocr.FloatRect, error) {
	var rects []ocr.FloatRect
	err := readObject(&rects, rectsFilePath)
	return rects, err
}

func SaveRects(rects []ocr.FloatRect) error {
	return writeObject(rects, rectsFilePath)
}

func SaveMenu(menu Menu) error {
	return writeObject(menu, menuFilePath)
}

func GetMenu() (Menu, error) {
	var menu Menu
	err := readObject(&menu, menuFilePath)
	return menu, err
}
