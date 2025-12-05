package go_new

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type rect struct {
	w, h float64
}

func TestNewVsPointer(t *testing.T) {
	r1 := rect{
		w: 1,
		h: 2,
	}

	r2 := new(rect)
	r2.w, r2.h = 3, 4

	r3 := &rect{}
	//&Type{}로 생성하면 초기값이 할당된 구조체의 포인터를 생성할 수 있음
	r4 := &rect{5, 6}

	assert.Equal(t, rect{1, 2}, r1)
	assert.Equal(t, rect{3, 4}, *r2)
	assert.Equal(t, rect{0, 0}, *r3)
	assert.Equal(t, rect{5, 6}, *r4)

}
