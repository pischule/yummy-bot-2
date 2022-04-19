package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestGetTextFromImage(t *testing.T) {
	bytes, err := ioutil.ReadFile("ocr_testdata/img.jpeg")
	if err != nil {
		t.Errorf("Error opening file: %v", err)
	}
	jsonBytes, _ := os.ReadFile("ocr_testdata/rois.json")
	var rois []FloatRect
	if err := json.Unmarshal(jsonBytes, &rois); err != nil {
		t.Error(err)
	}

	items := GetTextFromImage(bytes, rois)
	expected := []string{
		"салат из свеклы с сыром фета",
		"салат настроение",
		"салат с крабовыми палочками и кукурузой",
		"суп фасолевый с курицей",
		"овощи весенние на пару",
		"каша гречневая с морковью",
		"картофель запечённый под сырной корочкой",
		"бифштекс из птицы ся",
		"биточек из филе птицы с зеленью",
		"морской дуэт",
		"драники по-домашнему со сметаной",
		"филе птицы в сыре",
	}
	if !reflect.DeepEqual(items, expected) {
		t.Errorf("Expected %v, got %v", expected, items)
	}
}
