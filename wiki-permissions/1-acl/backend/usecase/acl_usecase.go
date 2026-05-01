package usecase

import (
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
)

// ACLUsecase는 ACL 관리(공유 추가/회수/조회) 흐름을 담당한다.
// 모든 메서드는 페이지 owner만 호출할 수 있도록 사전 검증을 수행한다.
type ACLUsecase struct {
	pages domain.PageRepository
	acls  domain.ACLRepository
}

// NewACLUsecase는 page/ACL 저장소를 주입받아 ACLUsecase를 생성한다.
func NewACLUsecase(pages domain.PageRepository, acls domain.ACLRepository) *ACLUsecase {
	return &ACLUsecase{pages: pages, acls: acls}
}

// checkOwner는 pageID의 페이지를 조회하고 requesterID가 owner인지 검증한다.
// owner가 아니면 ErrForbidden을 반환한다.
func (u *ACLUsecase) checkOwner(pageID, requesterID uint) (*domain.Page, error) {
	page, err := u.pages.FindByID(pageID)
	if err != nil {
		return nil, err
	}
	if page.OwnerID != requesterID {
		return nil, ErrForbidden
	}
	return page, nil
}

// List는 페이지 owner의 요청에 한해 해당 페이지의 모든 ACL 항목을 반환한다.
func (u *ACLUsecase) List(pageID, requesterID uint) ([]domain.ACLEntry, error) {
	if _, err := u.checkOwner(pageID, requesterID); err != nil {
		return nil, err
	}
	return u.acls.ListByPage(pageID)
}

// Grant는 페이지 owner의 요청에 한해 (page, targetUser, action) ACL을 추가한다.
// 잘못된 action 값은 ErrInvalidAction으로 거부한다 (입력 검증 → 권한 검증 → 실행 순).
func (u *ACLUsecase) Grant(pageID, requesterID, targetUserID uint, action domain.Action) error {
	if !action.Valid() {
		return ErrInvalidAction
	}
	if _, err := u.checkOwner(pageID, requesterID); err != nil {
		return err
	}
	return u.acls.Grant(pageID, targetUserID, action)
}

// Revoke는 페이지 owner의 요청에 한해 (page, targetUser)에 부여된 모든 action ACL을 삭제한다.
func (u *ACLUsecase) Revoke(pageID, requesterID, targetUserID uint) error {
	if _, err := u.checkOwner(pageID, requesterID); err != nil {
		return err
	}
	return u.acls.Revoke(pageID, targetUserID)
}
