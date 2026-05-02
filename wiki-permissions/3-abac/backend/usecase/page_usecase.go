package usecase

import (
	"errors"

	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/domain"
)

// ErrForbidden은 ABAC 정책이 read 액션을 거부했을 때 반환된다.
// HTTP 핸들러는 이 에러를 403 Forbidden으로 매핑한다.
var ErrForbidden = errors.New("forbidden")

// PageWithDecision은 페이지 + 호출자 기준 액션 결정을 함께 응답하는 구조다.
// 클라이언트는 CanRead/CanEdit의 reason을 그대로 사용자에게 표시할 수 있다 (ABAC의 미덕).
type PageWithDecision struct {
	Page    *domain.Page    `json:"page"`
	CanRead domain.Decision `json:"can_read"`
	CanEdit domain.Decision `json:"can_edit"`
}

// PageUsecase는 페이지 조회/수정 흐름을 담당하며, 각 액션 전에 ABAC 정책을 평가한다.
//
// 1편(ACL): EvaluateACL(page, userID, want, entries)
// 2편(RBAC): HasPermission(perms, want)
// 3편(ABAC): EvaluateABAC(user, page, action) — 사용자/페이지 속성 결합 평가, Decision 반환
type PageUsecase struct {
	pages domain.PageRepository
	users domain.UserRepository
}

func NewPageUsecase(pages domain.PageRepository, users domain.UserRepository) *PageUsecase {
	return &PageUsecase{pages: pages, users: users}
}

// List는 모든 페이지 중 호출자가 read 가능한 페이지만 반환한다.
//
// 1·2편의 List는 SQL JOIN/WHERE 또는 단순 List + permission lookup으로 처리됐지만,
// ABAC은 사용자 속성 + 페이지 속성을 결합 평가하므로 메모리 필터가 자연스럽다.
// 페이지 수가 많아지면 정책의 일부를 SQL로 변환해 사전 필터링하는 최적화가 가능하다.
func (u *PageUsecase) List(userID uint) ([]domain.Page, error) {
	user, err := u.users.FindByID(userID)
	if err != nil {
		return nil, err
	}
	all, err := u.pages.List()
	if err != nil {
		return nil, err
	}
	out := make([]domain.Page, 0, len(all))
	for i := range all {
		if domain.EvaluateABAC(user, &all[i], domain.ActionRead).Allowed {
			out = append(out, all[i])
		}
	}
	return out, nil
}

// Get은 read 정책 평가 후 페이지 + 양쪽 액션의 Decision을 반환한다.
// read가 거부되면 페이지 자체를 노출하지 않는다.
func (u *PageUsecase) Get(pageID, userID uint) (*PageWithDecision, error) {
	user, err := u.users.FindByID(userID)
	if err != nil {
		return nil, err
	}
	page, err := u.pages.FindByID(pageID)
	if err != nil {
		return nil, err
	}
	canRead := domain.EvaluateABAC(user, page, domain.ActionRead)
	if !canRead.Allowed {
		return nil, ErrForbidden
	}
	canEdit := domain.EvaluateABAC(user, page, domain.ActionEdit)
	return &PageWithDecision{Page: page, CanRead: canRead, CanEdit: canEdit}, nil
}

// Update는 edit 정책 평가 후 title/content를 갱신한다.
func (u *PageUsecase) Update(pageID, userID uint, title, content string) (*domain.Page, error) {
	user, err := u.users.FindByID(userID)
	if err != nil {
		return nil, err
	}
	page, err := u.pages.FindByID(pageID)
	if err != nil {
		return nil, err
	}
	if !domain.EvaluateABAC(user, page, domain.ActionEdit).Allowed {
		return nil, ErrForbidden
	}
	page.Title = title
	page.Content = content
	if err := u.pages.Update(page); err != nil {
		return nil, err
	}
	return page, nil
}
