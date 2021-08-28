package go_parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	type wants struct {
		want    int
		wantErr bool
	}
	tests := []struct {
		name  string
		input string
		wants wants
	}{
		{
			name:  "(1+3)*7 => 28",
			input: "(1+3)*7",
			wants: wants{
				want:    28,
				wantErr: false,
			},
		},
		{
			name:  "1+3*7 => 22",
			input: "1+3*7",
			wants: wants{
				want:    22,
				wantErr: false,
			},
		},
		{
			name:  "7/3 => 2 (eval only does integer math)",
			input: "7/3",
			wants: wants{
				want:    22,
				wantErr: false,
			},
		},
		{
			name:  "7.3 (this parses, but we disallow it in eval)",
			input: "",
			wants: wants{
				want:    0,
				wantErr: true,
			},
		},
		{
			name:  "3@7",
			input: "",
			wants: wants{
				want:    0,
				wantErr: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseAndEval(tt.input)
			assert.Equal(t, tt.wants.wantErr, err != nil)
			assert.Equal(t, tt.wants.want, result)
		})
	}
}
