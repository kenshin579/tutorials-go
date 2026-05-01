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
