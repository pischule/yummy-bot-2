package main

import (
	"image"
)

type FloatPoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type FloatRect struct {
	Min FloatPoint `json:"min"`
	Max FloatPoint `json:"max"`
}

func relativeToAbsolute(roi FloatRect, w int, h int) image.Rectangle {
	min := image.Point{X: int(roi.Min.X * float64(w)), Y: int(roi.Min.Y * float64(h))}
	max := image.Point{X: int(roi.Max.X * float64(w)), Y: int(roi.Max.Y * float64(h))}
	return image.Rectangle{Min: min, Max: max}
}
