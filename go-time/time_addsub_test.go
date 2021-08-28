package main

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	layoutHHMM = "15:04"
)

type TimeString string

func (ts TimeString) Add(d time.Duration) string {
	parse, _ := time.Parse(layoutHHMM, string(ts[:2]+":"+ts[2:]))
	result := parse.Add(d)
	return strings.Replace(result.Format(layoutHHMM), ":", "", 1)
}

func Test(t *testing.T) {
	type args struct {
		timeStr      string
		timeDuration time.Duration
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1000 + 15분 => 1015으로 출력되어야 한다",
			args: args{
				timeStr:      "1000",
				timeDuration: time.Minute * 15,
			},
			want: "1015",
		},
		{
			name: "1000 - 15분 => 0945으로 출력되어야 한다",
			args: args{
				timeStr:      "1000",
				timeDuration: -time.Minute * 15,
			},
			want: "0945",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, TimeString(tt.args.timeStr).Add(tt.args.timeDuration))
		})
	}
}
