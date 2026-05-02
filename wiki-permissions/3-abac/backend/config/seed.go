package config

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/domain"
	"github.com/kenshin579/tutorials-go/wiki-permissions/3-abac/backend/pkg/passwordhash"
)

// Seed는 ABAC 시리즈 3편 시연을 위한 시드 데이터를 삽입한다 (idempotent).
//
//	부서 2 (Engineering, Marketing)
//	사용자 4 — alice/bob (Eng/fulltime), carol (Mkt/fulltime), dave (Mkt/contract)
//	페이지 3 — Engineering Roadmap (internal/Eng), Q4 Marketing Plan (confidential/Mkt), Public Onboarding Guide (public)
//
// 4 사용자 × 3 페이지 = 12 케이스로 4개 ABAC 정책(owner/public/internal/confidential)이 모두 시연된다.
func Seed(db *gorm.DB) error {
	// 1) Departments
	depts := []domain.Department{{Name: "Engineering"}, {Name: "Marketing"}}
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&depts).Error; err != nil {
		return err
	}
	deptByName := map[string]uint{}
	for _, d := range depts {
		var found domain.Department
		if err := db.Where("name = ?", d.Name).First(&found).Error; err != nil {
			return err
		}
		deptByName[d.Name] = found.ID
	}

	// 2) Users
	hash, err := passwordhash.Hash("password")
	if err != nil {
		return err
	}
	users := []domain.User{
		{Email: "alice@example.com", Name: "Alice", PasswordHash: hash, DepartmentID: deptByName["Engineering"], EmploymentType: domain.EmploymentFulltime},
		{Email: "bob@example.com", Name: "Bob", PasswordHash: hash, DepartmentID: deptByName["Engineering"], EmploymentType: domain.EmploymentFulltime},
		{Email: "carol@example.com", Name: "Carol", PasswordHash: hash, DepartmentID: deptByName["Marketing"], EmploymentType: domain.EmploymentFulltime},
		{Email: "dave@example.com", Name: "Dave", PasswordHash: hash, DepartmentID: deptByName["Marketing"], EmploymentType: domain.EmploymentContract},
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

	// 3) Pages — title에 unique 없으므로 lookup-or-create
	engID, mktID := deptByName["Engineering"], deptByName["Marketing"]
	pageSpecs := []domain.Page{
		{Title: "Engineering Roadmap", Content: "2026 engineering plan", OwnerID: byEmail["alice@example.com"], Confidentiality: domain.ConfidentialityInternal, DepartmentID: &engID},
		{Title: "Q4 Marketing Plan", Content: "Q4 campaigns", OwnerID: byEmail["carol@example.com"], Confidentiality: domain.ConfidentialityConfidential, DepartmentID: &mktID},
		{Title: "Public Onboarding Guide", Content: "Welcome", OwnerID: byEmail["alice@example.com"], Confidentiality: domain.ConfidentialityPublic, DepartmentID: nil},
	}
	for _, p := range pageSpecs {
		var found domain.Page
		err := db.Where("title = ? AND owner_id = ?", p.Title, p.OwnerID).First(&found).Error
		switch {
		case err == nil:
			// 이미 존재
		case errors.Is(err, gorm.ErrRecordNotFound):
			created := p
			if err := db.Create(&created).Error; err != nil {
				return err
			}
		default:
			return err
		}
	}
	return nil
}
