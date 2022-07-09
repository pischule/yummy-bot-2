package ocr

import (
	"fmt"
	"image"
	"io"
	"regexp"
	"sort"
	"strings"
)

func getRectIndex(p image.Point, rects []image.Rectangle) int {
	for i, rect := range rects {
		if p.In(rect) {
			return i
		}
	}
	return -1
}

type BlockOfText struct {
	rectIndex int
	text      string
	y         int
}

func postProcessTextAbbyy(text string) []string {
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, ")", ")\n")
	lines := strings.Split(text, "\n")
	re := regexp.MustCompile(`^(\D*?)(\(.+)?$`)
	resultLines := make([]string, 0, 1)
	for _, line := range lines {
		submatch := re.FindStringSubmatch(line)
		if len(submatch) < 2 {
			continue
		}
		line = submatch[1]
		line = strings.Join(strings.Fields(line), " ")
		if line == "" {
			continue
		}
		resultLines = append(resultLines, line)
	}
	return resultLines
}

func extractLines(document abbyyDocument, relativeRects []FloatRect) []string {
	rects := make([]image.Rectangle, len(relativeRects))
	for i := 0; i < len(relativeRects); i++ {
		rects[i] = relativeToAbsolute(relativeRects[i], document.Page.Width, document.Page.Height)
	}
	page := document.Page
	blocks := make([]BlockOfText, 0, len(rects))

	for _, block := range page.Block {
		if block.BlockType != "Text" {
			continue
		}
		blockCenter := image.Point{
			X: (block.L + block.R) / 2,
			Y: (block.T + block.B) / 2,
		}
		rectIndex := getRectIndex(blockCenter, rects)
		if rectIndex == -1 {
			continue
		}
		previousParLineB := -100
		var sb strings.Builder
		for _, par := range block.Text.Par {
			if len(par.Line) == 0 {
				continue
			}
			firstLine := par.Line[0]
			// value from -1 to 11
			if firstLine.T-previousParLineB > 5 {
				sb.WriteString("\n")
			} else {
				sb.WriteString("")
			}
			for _, line := range par.Line {
				for _, cp := range line.Formatting.CharParams {
					charPoint := image.Point{
						X: cp.L,
						Y: cp.T,
					}
					if getRectIndex(charPoint, rects) == -1 {
						continue
					}
					if cp.Text == "" {
						sb.WriteString(" ")
					} else {
						sb.WriteString(cp.Text)
					}
				}
				sb.WriteString(" ")
				previousParLineB = line.B
			}
		}
		blocks = append(blocks, BlockOfText{
			rectIndex: rectIndex,
			text:      sb.String(),
			y:         blockCenter.Y,
		})
	}

	sort.Slice(blocks[:], func(i, j int) bool {
		if blocks[i].rectIndex != blocks[j].rectIndex {
			return blocks[i].rectIndex < blocks[j].rectIndex
		} else {
			return blocks[i].y < blocks[j].y
		}
	})

	allLines := make([]string, 0)
	for _, b := range blocks {
		for _, line := range postProcessTextAbbyy(b.text) {
			allLines = append(allLines, line)
		}
	}
	return allLines
}

func GetTextFromImageAbbyy(imageReader io.Reader, relativeRects []FloatRect, username string, password string) ([]string, error) {
	doc, err := recognizeFile(imageReader, username, password)
	if err != nil {
		fmt.Println("recognizeFile failed", err)
		return nil, err
	}
	return extractLines(doc, relativeRects), nil
}
