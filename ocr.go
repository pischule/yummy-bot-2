package main

import (
	"image"
	"sort"
	"strings"

	"github.com/pischule/gosseract/v2"
	"gocv.io/x/gocv"
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

func extractTextLinesFromImage(img gocv.Mat) []gocv.Mat {
	morph := gocv.NewMat()
	defer morph.Close()
	gocv.MorphologyExWithParams(img, &morph, gocv.MorphDilate,
		gocv.GetStructuringElement(gocv.MorphRect, image.Pt(10, 10)), 1, gocv.BorderIsolated)

	contours := gocv.FindContours(morph, gocv.RetrievalExternal, gocv.ChainApproxSimple)
	var boxes []image.Rectangle
	for i := 0; i < contours.Size(); i++ {
		rect := gocv.BoundingRect(contours.At(i))
		boxes = append(boxes, rect)
	}

	sort.Slice(boxes, func(i, j int) bool {
		return boxes[i].Min.Y < boxes[j].Min.Y
	})

	var lines []gocv.Mat
	for _, bbox := range boxes {
		line := img.Region(bbox)
		lines = append(lines, line)
	}

	return lines
}

func matToJpegBytes(mat gocv.Mat) []byte {
	buf, _ := gocv.IMEncode(gocv.JPEGFileExt, mat)
	return buf.GetBytes()
}

func postProcessText(text string) string {
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, "\n", " ")

	i := strings.Index(text, "(")
	if i != -1 {
		text = text[:i]
	}

	// filter out non-alphanumeric characters
	text = strings.Map(func(r rune) rune {
		if (r >= 'а' && r <= 'я') || (r >= 'А' && r <= 'Я') || r == ' ' || r == '-' || r == 'ё' || r == 'Ё' {
			return r
		}
		return ' '
	}, text)

	// remove duplicate spaces
	text = strings.Join(strings.Fields(text), " ")
	return text
}

func GetTextFromImage(bytes []byte, roiRects []FloatRect) []string {
	img, err := gocv.IMDecode(bytes, gocv.IMReadColor)
	if err != nil {
		panic(err)
	}
	gocv.CvtColor(img, &img, gocv.ColorBGRToGray)
	gocv.Threshold(img, &img, 150, 255, gocv.ThresholdBinary)

	absoluteRoiRectangles := make([]image.Rectangle, len(roiRects))
	for i, roi := range roiRects {
		absoluteRoiRectangles[i] = relativeToAbsolute(roi, img.Cols(), img.Rows())
	}

	roiImages := make([]gocv.Mat, len(absoluteRoiRectangles))
	for i, rectangle := range absoluteRoiRectangles {
		roiImages[i] = img.Region(rectangle)
	}

	var imgLines []gocv.Mat
	for _, roiImg := range roiImages {
		blockLines := extractTextLinesFromImage(roiImg)
		imgLines = append(imgLines, blockLines...)
	}

	var lines []string

	var ocr = gosseract.NewClient()
	defer ocr.Close()
	err = ocr.SetLanguage("rus")
	if err != nil {
		panic(err)
	}

	for _, imgLine := range imgLines {
		ocr.SetImageFromBytes(matToJpegBytes(imgLine))
		text, err := ocr.Text()
		if err != nil {
			continue
		}
		text = postProcessText(text)
		if len(text) > 2 {
			lines = append(lines, text)
		}
	}

	return lines
}
