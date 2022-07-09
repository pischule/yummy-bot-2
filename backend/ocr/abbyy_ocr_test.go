package ocr

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func Test_extractLines1(t *testing.T) {
	type args struct {
		filePath string
		url      string
	}
	var tests = []struct {
		name string
		args args
		want []string
	}{
		{
			name: "08072022",
			args: args{
				filePath: "testdata/08072022.xml",
				url:      "http://localhost/?r=124.416.319.107.508.415.310.108.138.574.296.96.511.573.302.100.145.713.327.111.517.712.292.114",
			},
			want: []string{"салат из свеклы с сыром фета",
				"салат «крабушка»",
				"салат «венгерский»",
				"бо^и^с^рикадельками",
				"картофель запеченный с прованскими травами",
				"рис с весенними овощами",
				"спагетти с маслом и зеленью",
				"рыбная котлета с зеленью",
				"ножка куриная фаршированная шампиньонами",
				"нагеттсы из филе птицы",
				"шницель «полесский»",
				"жаркое по-домашнему со свининой"},
		},
		{
			name: "07072022",
			args: args{
				filePath: "testdata/07072022.xml",
				url:      "https://pischule.github.io/yummy-bot-2/rects-tool/?r=119.330.327.114.500.332.296.110.126.515.286.131.502.517.294.134.129.698.314.111.508.696.285.114",
			},
			want: []string{
				"салат «бермудский»",
				"салат из свеклы с сыром фета",
				"сельдь под шубой",
				"борщ холодный",
				"картофельное пюре",
				"каша гречневая с морковью",
				"овощи «мехико» на пару",
				"бифштекс из птицы с яйцом",
				"котлета «папараць-кветка»",
				"куриные плечики запечённые в соусе",
				"филе птицы в сыре",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader, _ := os.Open(tt.args.filePath)
			bytes, _ := ioutil.ReadAll(reader)
			var doc abbyyDocument
			_ = xml.Unmarshal(bytes, &doc)
			relativeRects, _ := LoadRectsFromUri(tt.args.url)
			got := extractLines(doc, relativeRects)
			fmt.Println(strings.Join(got, "\n"))
			if len(got) != len(tt.want) {
				t.Errorf("extractLines() returned len %v, want %v", len(got), len(relativeRects))
				return
			}
			for i, item := range got {
				if tt.want[i] != strings.TrimSpace(item) {
					t.Errorf("extractLines() = %v, want %v", item, tt.want[i])
				}
			}
		})
	}
}
