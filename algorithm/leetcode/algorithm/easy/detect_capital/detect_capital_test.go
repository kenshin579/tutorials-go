package detect_capital

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_detectCapitalUse1(t *testing.T) {
	assert.True(t, DetectCapitalUse("USA"))
	assert.True(t, DetectCapitalUse("leetcode"))
	assert.True(t, DetectCapitalUse("Google"))
}

func Test_detectCapitalUse2(t *testing.T) {
	assert.False(t, DetectCapitalUse("FlagG"))
	assert.False(t, DetectCapitalUse("USa"))
	assert.False(t, DetectCapitalUse("usA"))
	assert.False(t, DetectCapitalUse("uSa"))
}

func Test_detectCapitalUse(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "True값인 경우",
			args: args{word: "USA"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DetectCapitalUse(tt.args.word); got != tt.want {
				t.Errorf("DetectCapitalUse() = %v, want %v", got, tt.want)
			}
		})
	}
}
