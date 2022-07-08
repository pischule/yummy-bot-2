package main

import (
	"testing"
)

func TestValidateName(t *testing.T) {
	var tests = []struct {
		input, wantName string
		wantErr         bool
	}{
		{"", "", true},
		{"   ", "", true},
		{"ё", "", true},
		{"vasya", "", true},
		{"Вася", "Вася", false},
	}

	for _, test := range tests {
		name, err := ValidateName(test.input)
		t.Run(test.input, func(t *testing.T) {
			if (err != nil) != test.wantErr {
				t.Errorf("got error %v, want error %v", err, test.wantErr)
			}
			if name != test.wantName {
				t.Errorf("got name %q, want name %q", name, test.wantName)
			}
		})
	}

}
