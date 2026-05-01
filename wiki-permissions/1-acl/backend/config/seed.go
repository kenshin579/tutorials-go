package config

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/domain"
	"github.com/kenshin579/tutorials-go/wiki-permissions/1-acl/backend/pkg/passwordhash"
)

// Seed는 ACL 시리즈 1편 시연을 위한 시드 데이터를 삽입한다.
// 사용자 4명(alice/bob/carol/dave, 비밀번호 모두 "password"), 페이지 3개,
// ACL entry 7개를 생성하며, 두 번째 호출 시에도 중복 없이 동일 결과를 보장한다(idempotent).
func Seed(db *gorm.DB) error {
	hash, err := passwordhash.Hash("password")
	if err != nil {
		return err
	}

	// User: email에 unique index가 있으므로 OnConflict DoNothing으로 idempotent 보장.
	users := []domain.User{
		{Email: "alice@example.com", Name: "Alice", PasswordHash: hash},
		{Email: "bob@example.com", Name: "Bob", PasswordHash: hash},
		{Email: "carol@example.com", Name: "Carol", PasswordHash: hash},
		{Email: "dave@example.com", Name: "Dave", PasswordHash: hash},
	}
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&users).Error; err != nil {
		return err
	}

	byEmail := map[string]uint{}
	for _, u := range users {
		var found domain.User
		if err := db.Where("email = ?", u.Email).First(&found).Error; err != nil {
			return err
		}
		byEmail[u.Email] = found.ID
	}

	// Page: title에 unique index가 없어 OnConflict로 idempotent 못 함.
	// "title + owner_id 조합으로 조회 후 없으면 삽입" 패턴 사용.
	pageSpecs := []domain.Page{
		{Title: "Engineering Roadmap", Content: "2026 engineering plan", OwnerID: byEmail["alice@example.com"]},
		{Title: "Q4 Marketing Plan", Content: "Q4 campaigns", OwnerID: byEmail["carol@example.com"]},
		{Title: "Public Onboarding Guide", Content: "Welcome", OwnerID: byEmail["alice@example.com"]},
	}
	byTitle := map[string]uint{}
	for _, p := range pageSpecs {
		var found domain.Page
		err := db.Where("title = ? AND owner_id = ?", p.Title, p.OwnerID).First(&found).Error
		switch {
		case err == nil:
			byTitle[p.Title] = found.ID
		case errors.Is(err, gorm.ErrRecordNotFound):
			created := p
			if err := db.Create(&created).Error; err != nil {
				return err
			}
			byTitle[p.Title] = created.ID
		default:
			return err
		}
	}

	// ACLEntry: idx_page_user_action 복합 unique 인덱스 → OnConflict DoNothing으로 idempotent.
	type aclSpec struct {
		page string
		user string
		act  domain.Action
	}
	specs := []aclSpec{
		{"Engineering Roadmap", "bob@example.com", domain.ActionEdit},
		{"Engineering Roadmap", "carol@example.com", domain.ActionRead},
		{"Q4 Marketing Plan", "alice@example.com", domain.ActionRead},
		{"Q4 Marketing Plan", "bob@example.com", domain.ActionRead},
		{"Public Onboarding Guide", "bob@example.com", domain.ActionRead},
		{"Public Onboarding Guide", "carol@example.com", domain.ActionRead},
		{"Public Onboarding Guide", "dave@example.com", domain.ActionRead},
	}
	entries := make([]domain.ACLEntry, 0, len(specs))
	for _, s := range specs {
		entries = append(entries, domain.ACLEntry{
			PageID: byTitle[s.page],
			UserID: byEmail[s.user],
			Action: s.act,
		})
	}
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&entries).Error
}
