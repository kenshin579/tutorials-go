package go1_25

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NilPointer_잘못된_패턴_panic_발생(t *testing.T) {
	// Go 1.25에서 수정: nil 포인터 사용 시 올바르게 panic 발생
	assert.Panics(t, func() {
		f, _ := os.Open("존재하지않는파일.txt")
		// 에러 확인 없이 nil 포인터의 메서드 호출 → panic!
		_ = f.Name()
	}, "nil 포인터 메서드 호출 시 panic이 발생해야 한다")
}

func Test_NilPointer_올바른_패턴(t *testing.T) {
	// 올바른 패턴: 에러를 먼저 확인한 후 포인터 사용
	f, err := os.Open("존재하지않는파일.txt")
	if err != nil {
		t.Logf("예상된 에러: %v", err)
		return // 에러 발생 시 조기 반환
	}
	defer f.Close()

	// 에러가 없을 때만 포인터 사용
	name := f.Name()
	assert.NotEmpty(t, name)
}
