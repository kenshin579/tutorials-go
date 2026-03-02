package go_testing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAverage_Basic - 기본적인 테이블 기반 테스트 (t.Run 없이)
func TestAverage_Basic(t *testing.T) {
	for _, tt := range []struct {
		Nos    []int
		Result int
	}{
		{Nos: []int{2, 4}, Result: 3},
		{Nos: []int{1, 2, 5}, Result: 2},
		{Nos: []int{1}, Result: 1},
		{Nos: []int{}, Result: 0},
		{Nos: []int{2, -2}, Result: 0},
	} {
		if avg := Average(tt.Nos...); avg != tt.Result {
			t.Fatalf("expected average of %v to be %d, got %d\n", tt.Nos, tt.Result, avg)
		}
	}
}

// TestAverage_TableDriven - t.Run()을 사용한 서브테스트 패턴
// 각 테스트 케이스에 이름을 부여하여 개별 실행 및 디버깅이 용이하다.
// 실행: go test -run TestAverage_TableDriven/빈_슬라이스 -v
func TestAverage_TableDriven(t *testing.T) {
	tests := []struct {
		name   string
		nos    []int
		want   int
	}{
		{name: "두_수의_평균", nos: []int{2, 4}, want: 3},
		{name: "세_수의_평균", nos: []int{1, 2, 5}, want: 2},
		{name: "단일_값", nos: []int{1}, want: 1},
		{name: "빈_슬라이스", nos: []int{}, want: 0},
		{name: "합이_0인_경우", nos: []int{2, -2}, want: 0},
		{name: "큰_수", nos: []int{100, 200, 300}, want: 200},
		{name: "음수_포함", nos: []int{-10, -20, -30}, want: -20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Average(tt.nos...)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestAverage_Parallel - t.Parallel()을 사용한 병렬 테스트
// 각 서브테스트가 독립적으로 병렬 실행된다.
func TestAverage_Parallel(t *testing.T) {
	tests := []struct {
		name string
		nos  []int
		want int
	}{
		{name: "두_수의_평균", nos: []int{2, 4}, want: 3},
		{name: "세_수의_평균", nos: []int{1, 2, 5}, want: 2},
		{name: "단일_값", nos: []int{1}, want: 1},
		{name: "빈_슬라이스", nos: []int{}, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := Average(tt.nos...)
			assert.Equal(t, tt.want, got)
		})
	}
}
