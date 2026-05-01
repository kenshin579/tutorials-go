package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/config"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/repository"
)

func setupPageEnv(t *testing.T) (*PageUsecase, map[string]uint) {
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	require.NoError(t, config.Seed(db))

	pages := repository.NewPageRepository(db)
	acls := repository.NewACLRepository(db)
	uc := NewPageUsecase(pages, acls)

	users := repository.NewUserRepository(db)

	emailToID := map[string]uint{}
	for _, e := range []string{"alice@example.com", "bob@example.com", "carol@example.com", "dave@example.com"} {
		u, _ := users.FindByEmail(e)
		emailToID[e] = u.ID
	}
	titleToID := map[string]uint{}
	for _, title := range []string{"Engineering Roadmap", "Q4 Marketing Plan", "Public Onboarding Guide"} {
		var p domain.Page
		require.NoError(t, db.Where("title = ?", title).First(&p).Error)
		titleToID[title] = p.ID
	}
	return uc, mergeIDs(emailToID, titleToID)
}

func mergeIDs(a, b map[string]uint) map[string]uint {
	out := map[string]uint{}
	for k, v := range a {
		out["user:"+k] = v
	}
	for k, v := range b {
		out["page:"+k] = v
	}
	return out
}

func TestPageUsecase_Get_OwnerCanRead(t *testing.T) {
	uc, ids := setupPageEnv(t)
	p, err := uc.Get(ids["page:Engineering Roadmap"], ids["user:alice@example.com"])
	require.NoError(t, err)
	assert.Equal(t, "Engineering Roadmap", p.Title)
}

func TestPageUsecase_Get_GrantedUserCanRead(t *testing.T) {
	uc, ids := setupPageEnv(t)
	p, err := uc.Get(ids["page:Engineering Roadmap"], ids["user:carol@example.com"])
	require.NoError(t, err)
	assert.NotNil(t, p)
}

func TestPageUsecase_Get_NoGrantReturnsForbidden(t *testing.T) {
	uc, ids := setupPageEnv(t)
	_, err := uc.Get(ids["page:Engineering Roadmap"], ids["user:dave@example.com"])
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestPageUsecase_Update_RequiresEdit(t *testing.T) {
	uc, ids := setupPageEnv(t)
	// carol은 EngRoadmap에 read만 있음 → edit 시도 시 403
	_, err := uc.Update(ids["page:Engineering Roadmap"], ids["user:carol@example.com"], "new title", "new content")
	assert.ErrorIs(t, err, ErrForbidden)

	// bob은 edit 권한 있음
	p, err := uc.Update(ids["page:Engineering Roadmap"], ids["user:bob@example.com"], "Updated", "x")
	require.NoError(t, err)
	assert.Equal(t, "Updated", p.Title)
}
