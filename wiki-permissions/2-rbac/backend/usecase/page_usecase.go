package usecase

import (
	"errors"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/domain"
)

// ErrForbidden은 사용자가 해당 액션의 RBAC 권한을 가지지 않을 때 반환된다.
// HTTP 핸들러는 이 에러를 403 Forbidden으로 매핑한다.
var ErrForbidden = errors.New("forbidden")

// HasPermission은 사용자의 효과적 권한 집합에 want("resource:action") 키가 포함되는지 확인한다.
//
// 1편(ACL)의 EvaluateACL은 owner short-circuit + edit→read 함의 등 5가지 규칙의 30줄 함수였지만,
// RBAC에서는 사용자→role→permission JOIN을 PermissionRepository가 끝내므로 이 단계는 한 줄 lookup이다.
// 비교의 무게가 "평가 함수"에서 "데이터 모델"로 옮겨가는 것이 시리즈의 핵심 메시지.
func HasPermission(perms []domain.Permission, want string) bool {
	for _, p := range perms {
		if p.Key() == want {
			return true
		}
	}
	return false
}

// PageUsecase는 페이지 CRUD 흐름을 담당하며, 모든 액션 전에 RBAC 권한을 평가한다.
// 1편 PageUsecase와의 차이: ACL 평가가 RBAC 평가로 교체됨, owner_id 무시, Create/Delete 추가.
type PageUsecase struct {
	pages domain.PageRepository
	perms domain.PermissionRepository
}

// NewPageUsecase는 page/permission 저장소를 주입받아 PageUsecase를 생성한다.
func NewPageUsecase(pages domain.PageRepository, perms domain.PermissionRepository) *PageUsecase {
	return &PageUsecase{pages: pages, perms: perms}
}

// requirePerm은 사용자의 효과적 권한을 조회한 뒤 want가 포함되는지 확인한다.
// 미포함 시 ErrForbidden을 반환한다.
func (u *PageUsecase) requirePerm(userID uint, want string) error {
	ps, err := u.perms.FindByUserID(userID)
	if err != nil {
		return err
	}
	if !HasPermission(ps, want) {
		return ErrForbidden
	}
	return nil
}

// List는 pages:read 권한이 있으면 모든 페이지를 반환한다.
// (1편에서는 본인 owner / ACL 매칭 페이지만 반환했지만, RBAC은 role 기반이라 모두 노출.)
func (u *PageUsecase) List(userID uint) ([]domain.Page, error) {
	if err := u.requirePerm(userID, "pages:read"); err != nil {
		return nil, err
	}
	return u.pages.List()
}

// Get은 pages:read 권한이 있으면 페이지 상세를 반환한다.
func (u *PageUsecase) Get(pageID, userID uint) (*domain.Page, error) {
	if err := u.requirePerm(userID, "pages:read"); err != nil {
		return nil, err
	}
	return u.pages.FindByID(pageID)
}

// Create는 pages:create 권한이 있으면 새 페이지를 생성한다 (owner = 생성자).
func (u *PageUsecase) Create(userID uint, title, content string) (*domain.Page, error) {
	if err := u.requirePerm(userID, "pages:create"); err != nil {
		return nil, err
	}
	p := &domain.Page{Title: title, Content: content, OwnerID: userID}
	if err := u.pages.Create(p); err != nil {
		return nil, err
	}
	return p, nil
}

// Update는 pages:edit 권한이 있으면 페이지를 갱신한다.
//
// 한계 — RBAC만으로는 "내가 만든 페이지만 수정"을 표현할 수 없다.
// editor role을 가진 두 사용자가 서로의 페이지를 수정할 수 있다 (3편 ABAC에서 owner 속성으로 해결).
func (u *PageUsecase) Update(pageID, userID uint, title, content string) (*domain.Page, error) {
	if err := u.requirePerm(userID, "pages:edit"); err != nil {
		return nil, err
	}
	page, err := u.pages.FindByID(pageID)
	if err != nil {
		return nil, err
	}
	page.Title = title
	page.Content = content
	if err := u.pages.Update(page); err != nil {
		return nil, err
	}
	return page, nil
}

// Delete는 pages:delete 권한이 있으면 페이지를 삭제한다 (admin 전용).
func (u *PageUsecase) Delete(pageID, userID uint) error {
	if err := u.requirePerm(userID, "pages:delete"); err != nil {
		return err
	}
	return u.pages.Delete(pageID)
}
