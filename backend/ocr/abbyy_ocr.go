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
	text = strings.ReplaceAll(text, "^", "")
	lines := strings.Split(text, "\n")
	re := regexp.MustCompile(`^(.*?)(\(.+)?$`)
	resultLines := make([]string, 0, 1)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		line = re.ReplaceAllString(line, `$1`)
		if line == "" {
			continue
		}
		resultLines = append(resultLines, line)
	}
	return resultLines
}

func extractLines(document abbyyDocument, rects []image.Rectangle) []string {
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
		var sb strings.Builder
		for _, par := range block.Text.Par {
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
			}
			sb.WriteString("\n")
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

	absoluteRects := make([]image.Rectangle, len(relativeRects))
	for i := 0; i < len(relativeRects); i++ {
		absoluteRects[i] = relativeToAbsolute(relativeRects[i], doc.Page.Width, doc.Page.Height)
	}

	return extractLines(doc, absoluteRects), nil
}
