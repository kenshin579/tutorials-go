package go1_25

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GOMAXPROCS_현재값_조회(t *testing.T) {
	// GOMAXPROCS(0)은 현재 값을 변경하지 않고 반환
	current := runtime.GOMAXPROCS(0)
	assert.Positive(t, current, "GOMAXPROCS는 양수여야 한다")
	t.Logf("현재 GOMAXPROCS: %d", current)
}

func Test_SetDefaultGOMAXPROCS_기본값_복원(t *testing.T) {
	original := runtime.GOMAXPROCS(0)

	// GOMAXPROCS를 수동으로 변경
	runtime.GOMAXPROCS(2)
	assert.Equal(t, 2, runtime.GOMAXPROCS(0))

	// SetDefaultGOMAXPROCS()로 기본값(CPU 수 기반)으로 복원
	runtime.SetDefaultGOMAXPROCS()

	restored := runtime.GOMAXPROCS(0)
	t.Logf("복원된 GOMAXPROCS: %d (원래: %d)", restored, original)
	assert.Positive(t, restored, "복원된 값은 양수여야 한다")

	// 원래 값으로 되돌리기
	runtime.GOMAXPROCS(original)
}
