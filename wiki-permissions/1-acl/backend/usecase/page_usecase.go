package usecase

import (
	"errors"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
)

// ErrForbidden은 사용자가 해당 리소스에 대한 권한을 가지지 않을 때 반환된다.
// HTTP 핸들러는 이 에러를 403 Forbidden으로 매핑한다.
var ErrForbidden = errors.New("forbidden")

// ErrInvalidAction은 클라이언트가 알려지지 않은 ACL action 값을 보냈을 때 반환된다.
// HTTP 핸들러는 이 에러를 400 Bad Request로 매핑한다 (도메인의 ErrNotFound와 분리되는 입력 검증 실패다).
var ErrInvalidAction = errors.New("invalid action")

// PageUsecase는 페이지 조회/수정/목록 조회 흐름을 담당하며, 모든 액션 전에 ACL을 평가한다.
type PageUsecase struct {
	pages domain.PageRepository
	acls  domain.ACLRepository
}

// NewPageUsecase는 page/ACL 저장소를 주입받아 PageUsecase를 생성한다.
func NewPageUsecase(pages domain.PageRepository, acls domain.ACLRepository) *PageUsecase {
	return &PageUsecase{pages: pages, acls: acls}
}

// ListAccessible은 userID가 owner이거나 ACL을 받은 모든 페이지 목록을 반환한다.
// 별도 ACL 평가가 필요 없다(repository 쿼리 자체가 owner OR ACL 매칭).
func (u *PageUsecase) ListAccessible(userID uint) ([]domain.Page, error) {
	return u.pages.ListAccessibleBy(userID)
}

// Get은 pageID/userID에 대해 read 권한을 평가한 뒤 페이지 상세를 반환한다.
// 권한이 없으면 ErrForbidden을 반환한다.
func (u *PageUsecase) Get(pageID, userID uint) (*domain.Page, error) {
	page, err := u.pages.FindByID(pageID)
	if err != nil {
		return nil, err
	}
	entries, err := u.acls.FindByPageAndUser(pageID, userID)
	if err != nil {
		return nil, err
	}
	if !domain.EvaluateACL(page, userID, domain.ActionRead, entries) {
		return nil, ErrForbidden
	}
	return page, nil
}

// Update은 edit 권한을 평가한 뒤 title/content를 갱신하고 갱신된 페이지를 반환한다.
// 권한이 없으면 ErrForbidden을 반환한다.
func (u *PageUsecase) Update(pageID, userID uint, title, content string) (*domain.Page, error) {
	page, err := u.pages.FindByID(pageID)
	if err != nil {
		return nil, err
	}
	entries, err := u.acls.FindByPageAndUser(pageID, userID)
	if err != nil {
		return nil, err
	}
	if !domain.EvaluateACL(page, userID, domain.ActionEdit, entries) {
		return nil, ErrForbidden
	}
	page.Title = title
	page.Content = content
	if err := u.pages.Update(page); err != nil {
		return nil, err
	}
	return page, nil
}
