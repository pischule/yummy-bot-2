package ocr

import (
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
)

func LoadRectsFromUri(uri string) ([]FloatRect, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return nil, err
	}
	rectsString := m["r"]
	if rectsString == nil || len(rectsString) != 1 {
		return nil, fmt.Errorf("no query param rects found")
	}
	stringNumbers := strings.Split(rectsString[0], ".")
	if len(stringNumbers)%4 != 0 {
		return nil, fmt.Errorf("int numbers not divisable by 4")
	}
	numbers := make([]float64, 0, len(stringNumbers))
	for _, n := range stringNumbers {
		nFloat, err := strconv.ParseFloat(n, 64)
		if err != nil {
			return nil, fmt.Errorf("error while parsing number %s", n)
		}
		numbers = append(numbers, nFloat/1000)
	}

	rects := make([]FloatRect, 0, len(numbers)/4)

	for i := 0; i < len(numbers); i += 4 {
		rects = append(rects, FloatRect{
			Min: FloatPoint{
				X: numbers[i],
				Y: numbers[i+1],
			},
			Max: FloatPoint{
				X: numbers[i] + numbers[i+2],
				Y: numbers[i+1] + numbers[i+3],
			},
		})
	}
	return rects, nil
}

func RectsToUri(rects []FloatRect) string {
	points := make([]string, 0, len(rects)*4)
	for _, r := range rects {
		points = append(points, fmt.Sprint(math.Round(r.Min.X*1000)))
		points = append(points, fmt.Sprint(math.Round(r.Min.Y*1000)))
		points = append(points, fmt.Sprint(math.Round((r.Max.X-r.Min.X)*1000)))
		points = append(points, fmt.Sprint(math.Round((r.Max.Y-r.Min.Y)*1000)))
	}
	return "https://pischule.github.io/yummy-bot-2/?r=" + strings.Join(points, ".")
}
