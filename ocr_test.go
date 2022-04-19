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
		"салат настроение",
		"салат из свеклы с сыром фета",
		"салат с крабовыми палочками и кукурузой",
		"суп фасолевый с курицей",
		"картофель запечённый под сырной корочкой",
		"каша гречневая с морковью",
		"овощи весенние на пару",
		"морской дуэт",
		"биточек из филе птицы с зеленью",
		"бифштекс из птицы ся",
		"филе птицы в сыре",
		"драники по-домашнему со сметаной",
	}
	if !reflect.DeepEqual(items, expected) {
		t.Errorf("Expected %v, got %v", expected, items)
	}
}
