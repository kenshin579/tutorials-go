package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluateACL_OwnerHasFullAccess(t *testing.T) {
	page := &Page{ID: 1, OwnerID: 100}
	assert.True(t, EvaluateACL(page, 100, ActionEdit, nil))
	assert.True(t, EvaluateACL(page, 100, ActionRead, nil))
}

func TestEvaluateACL_ExplicitGrant(t *testing.T) {
	page := &Page{ID: 1, OwnerID: 100}
	entries := []ACLEntry{{PageID: 1, UserID: 200, Action: ActionRead}}
	assert.True(t, EvaluateACL(page, 200, ActionRead, entries))
	assert.False(t, EvaluateACL(page, 200, ActionEdit, entries))
}

func TestEvaluateACL_EditImpliesRead(t *testing.T) {
	page := &Page{ID: 1, OwnerID: 100}
	entries := []ACLEntry{{PageID: 1, UserID: 200, Action: ActionEdit}}
	assert.True(t, EvaluateACL(page, 200, ActionRead, entries))
	assert.True(t, EvaluateACL(page, 200, ActionEdit, entries))
}

func TestEvaluateACL_NoGrant(t *testing.T) {
	page := &Page{ID: 1, OwnerID: 100}
	assert.False(t, EvaluateACL(page, 999, ActionRead, nil))
}

// 다른 페이지/다른 사용자에 대한 entry는 평가에 영향이 없어야 한다 (방어적 필터 보장).
func TestEvaluateACL_IgnoresEntriesForOtherPagesOrUsers(t *testing.T) {
	page := &Page{ID: 1, OwnerID: 100}
	entries := []ACLEntry{
		{PageID: 2, UserID: 200, Action: ActionRead}, // 다른 페이지
		{PageID: 1, UserID: 999, Action: ActionRead}, // 다른 사용자
	}
	assert.False(t, EvaluateACL(page, 200, ActionRead, entries))
}

func TestEvaluateACL_NilPage(t *testing.T) {
	assert.False(t, EvaluateACL(nil, 100, ActionRead, nil))
}
