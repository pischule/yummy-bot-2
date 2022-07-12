package ocr

import (
	"reflect"
	"testing"
)

func Test_loadRectsFromUri(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name    string
		args    args
		want    []FloatRect
		wantErr bool
	}{
		{
			name: "correct url with 2 points",
			args: args{uri: "http://localhost:3000/?r=129.421.319.100.509.417.310.105"},
			want: []FloatRect{
				{
					Min: FloatPoint{
						X: 0.129,
						Y: 0.421,
					},
					Max: FloatPoint{
						X: 0.448,
						Y: 0.521,
					},
				},
				{
					Min: FloatPoint{
						X: 0.509,
						Y: 0.417,
					},
					Max: FloatPoint{
						X: 0.819,
						Y: 0.522,
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "incorrect url with 5 numbers",
			args:    args{"https://pischule.github.io/yummy-bot-2/rects-tool/?r=11.414.314.107.508"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "incorrect url with 5 numbers",
			args:    args{"https://pischule.github.io/yummy-bot-2/rects-tool/?r="},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty uri",
			args:    args{""},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "not uri",
			args:    args{"hello world"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "uri without query param",
			args:    args{"https://pischule.github.io/yummy-bot-2/rects-tool/"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadRectsFromUri(tt.args.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadRectsFromUri() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadRectsFromUri() got = %v, want %v", got, tt.want)
			}
		})
	}

}

func Test_rectsToUri(t *testing.T) {
	type args struct {
		rects []FloatRect
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "export 2 points",
			args: args{
				rects: []FloatRect{
					{
						Min: FloatPoint{
							X: 0.129,
							Y: 0.421,
						},
						Max: FloatPoint{
							X: 0.448,
							Y: 0.521,
						},
					},
					{
						Min: FloatPoint{
							X: 0.509,
							Y: 0.417,
						},
						Max: FloatPoint{
							X: 0.819,
							Y: 0.522,
						},
					},
				},
			},
			want: "https://pischule.github.io/yummy-bot-2/?r=129.421.319.100.509.417.310.105",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RectsToUri(tt.args.rects); got != tt.want {
				t.Errorf("rectsToUri() = %v, want %v", got, tt.want)
			}
		})
	}
}
