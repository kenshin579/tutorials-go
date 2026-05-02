package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/config"
	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/domain"
	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/repository"
)

func setupPageEnv(t *testing.T) (uc *PageUsecase, ids map[string]uint) {
	t.Helper()
	db, err := config.OpenDB(":memory:")
	require.NoError(t, err)
	require.NoError(t, config.Seed(db))

	users := repository.NewUserRepository(db)
	pages := repository.NewPageRepository(db)
	uc = NewPageUsecase(pages, users)

	ids = map[string]uint{}
	for _, e := range []string{"alice@example.com", "bob@example.com", "carol@example.com", "dave@example.com"} {
		u, err := users.FindByEmail(e)
		require.NoError(t, err)
		ids[e] = u.ID
	}
	for _, title := range []string{"Engineering Roadmap", "Q4 Marketing Plan", "Public Onboarding Guide"} {
		var p domain.Page
		require.NoError(t, db.Where("title = ?", title).First(&p).Error)
		ids[title] = p.ID
	}
	return uc, ids
}

// 4 사용자 × 3 페이지 = 12 케이스 매트릭스.
// alice (Eng/fulltime, owner of EngRoadmap & Onboarding)
// bob   (Eng/fulltime)
// carol (Mkt/fulltime, owner of Q4)
// dave  (Mkt/contract)
//
//	페이지: EngRoadmap (internal/Eng), Q4 (confidential/Mkt), Onboarding (public)
//
// 기대값:
//
//	alice 2 (EngRoadmap owner + Onboarding public; Q4는 다른 부서 confidential)
//	bob   2 (EngRoadmap 같은 부서 internal + Onboarding public; Q4는 다른 부서)
//	carol 2 (Q4 owner + Onboarding public; EngRoadmap는 다른 부서)
//	dave  1 (Onboarding public; EngRoadmap 다른 부서, Q4 contract)
func TestPageUsecase_List_FiltersByPolicy(t *testing.T) {
	uc, ids := setupPageEnv(t)

	cases := []struct {
		email string
		count int
	}{
		{"alice@example.com", 2},
		{"bob@example.com", 2},
		{"carol@example.com", 2},
		{"dave@example.com", 1},
	}
	for _, c := range cases {
		got, err := uc.List(ids[c.email])
		require.NoError(t, err)
		assert.Len(t, got, c.count, "user %s should see %d pages", c.email, c.count)
	}
}

func TestPageUsecase_Get_DecisionWithCanEdit(t *testing.T) {
	uc, ids := setupPageEnv(t)

	// bob (Eng/fulltime)이 EngRoadmap (internal/Eng)
	res, err := uc.Get(ids["Engineering Roadmap"], ids["bob@example.com"])
	require.NoError(t, err)
	assert.True(t, res.CanRead.Allowed)
	assert.True(t, res.CanEdit.Allowed)
	assert.Equal(t, "internal", res.CanRead.Policy)

	// alice가 Public Onboarding (owner) → owner 정책으로 read+edit 모두 가능
	res2, err := uc.Get(ids["Public Onboarding Guide"], ids["alice@example.com"])
	require.NoError(t, err)
	assert.Equal(t, "owner", res2.CanRead.Policy)
	assert.True(t, res2.CanEdit.Allowed)

	// dave가 Onboarding (public, 다른 사람 owner) → public read 가능, edit 거부
	res3, err := uc.Get(ids["Public Onboarding Guide"], ids["dave@example.com"])
	require.NoError(t, err)
	assert.True(t, res3.CanRead.Allowed)
	assert.False(t, res3.CanEdit.Allowed)
	assert.Equal(t, "public", res3.CanRead.Policy)
}

func TestPageUsecase_Get_ForbiddenForOtherDepartment(t *testing.T) {
	uc, ids := setupPageEnv(t)
	// bob (Eng)이 Q4 (confidential/Mkt) — 다른 부서
	_, err := uc.Get(ids["Q4 Marketing Plan"], ids["bob@example.com"])
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestPageUsecase_Update_RequiresEditDecision(t *testing.T) {
	uc, ids := setupPageEnv(t)

	// alice (owner)가 자기 EngRoadmap 갱신 OK
	p, err := uc.Update(ids["Engineering Roadmap"], ids["alice@example.com"], "Updated", "x")
	require.NoError(t, err)
	assert.Equal(t, "Updated", p.Title)

	// bob (Eng)이 Q4 (Mkt) 갱신 시도 → 다른 부서 거부
	_, err = uc.Update(ids["Q4 Marketing Plan"], ids["bob@example.com"], "x", "y")
	assert.ErrorIs(t, err, ErrForbidden)

	// dave (Mkt/contract)가 Q4 갱신 시도 → contract → confidential 거부
	_, err = uc.Update(ids["Q4 Marketing Plan"], ids["dave@example.com"], "x", "y")
	assert.ErrorIs(t, err, ErrForbidden)
}
