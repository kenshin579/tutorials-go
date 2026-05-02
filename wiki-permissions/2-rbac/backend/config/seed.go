package config

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/domain"
	"github.com/kenshin579/tutorials-go/wiki-permissions/2-rbac/backend/pkg/passwordhash"
)

// Seed는 RBAC 시리즈 2편 시연을 위한 시드 데이터를 삽입한다 (idempotent).
//
//	사용자 4 (alice/bob/carol/dave, 비밀번호 모두 "password")
//	페이지 3 (Engineering Roadmap, Q4 Marketing Plan, Public Onboarding Guide)
//	permissions 6 (pages:read/create/edit/delete, users:read/manage)
//	roles 3 (admin/editor/viewer)
//	role-permission 매트릭스 + user-role 매핑
//
// 두 번째 호출 시에도 중복 row 없이 동일 결과를 보장한다 (앱 부팅 시 단일 호출 가정).
func Seed(db *gorm.DB) error {
	hash, err := passwordhash.Hash("password")
	if err != nil {
		return err
	}

	// 1) Users — Email uniqueIndex + OnConflict
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

	// 2) Pages — title에 unique 없으므로 lookup-or-create
	pageSpecs := []domain.Page{
		{Title: "Engineering Roadmap", Content: "2026 engineering plan", OwnerID: byEmail["alice@example.com"]},
		{Title: "Q4 Marketing Plan", Content: "Q4 campaigns", OwnerID: byEmail["carol@example.com"]},
		{Title: "Public Onboarding Guide", Content: "Welcome", OwnerID: byEmail["alice@example.com"]},
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

	// 3) Permissions — (resource, action) 복합 unique + OnConflict
	permSpecs := []domain.Permission{
		{Resource: "pages", Action: "read"},
		{Resource: "pages", Action: "create"},
		{Resource: "pages", Action: "edit"},
		{Resource: "pages", Action: "delete"},
		{Resource: "users", Action: "read"},
		{Resource: "users", Action: "manage"},
	}
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&permSpecs).Error; err != nil {
		return err
	}
	permByKey := map[string]uint{}
	for _, p := range permSpecs {
		var found domain.Permission
		if err := db.Where("resource = ? AND action = ?", p.Resource, p.Action).First(&found).Error; err != nil {
			return err
		}
		permByKey[p.Key()] = found.ID
	}

	// 4) Roles — name uniqueIndex + OnConflict
	roleSpecs := []domain.Role{
		{Name: "admin", Description: "All permissions"},
		{Name: "editor", Description: "Pages CRUD except delete"},
		{Name: "viewer", Description: "Read pages only"},
	}
	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&roleSpecs).Error; err != nil {
		return err
	}
	roleByName := map[string]uint{}
	for _, r := range roleSpecs {
		var found domain.Role
		if err := db.Where("name = ?", r.Name).First(&found).Error; err != nil {
			return err
		}
		roleByName[r.Name] = found.ID
	}

	// 5) RolePermission — 권한 매트릭스. GORM Association.Append는 중복 시 row 추가하므로
	//    이미 매핑된 경우는 skip해야 idempotent.
	matrix := map[string][]string{
		"admin":  {"pages:read", "pages:create", "pages:edit", "pages:delete", "users:read", "users:manage"},
		"editor": {"pages:read", "pages:create", "pages:edit"},
		"viewer": {"pages:read"},
	}
	for roleName, keys := range matrix {
		role := &domain.Role{ID: roleByName[roleName]}
		var existing []domain.Permission
		if err := db.Model(role).Association("Permissions").Find(&existing); err != nil {
			return err
		}
		existingIDs := map[uint]struct{}{}
		for _, e := range existing {
			existingIDs[e.ID] = struct{}{}
		}
		for _, k := range keys {
			pid := permByKey[k]
			if _, ok := existingIDs[pid]; ok {
				continue
			}
			perm := &domain.Permission{ID: pid}
			if err := db.Model(role).Association("Permissions").Append(perm); err != nil {
				return err
			}
		}
	}

	// 6) UserRole 매트릭스 — 위와 같은 패턴.
	userRoles := map[string]string{
		"alice@example.com": "admin",
		"bob@example.com":   "editor",
		"carol@example.com": "viewer",
		"dave@example.com":  "viewer",
	}
	for email, roleName := range userRoles {
		user := &domain.User{ID: byEmail[email]}
		var existing []domain.Role
		if err := db.Model(user).Association("Roles").Find(&existing); err != nil {
			return err
		}
		existingIDs := map[uint]struct{}{}
		for _, e := range existing {
			existingIDs[e.ID] = struct{}{}
		}
		rid := roleByName[roleName]
		if _, ok := existingIDs[rid]; ok {
			continue
		}
		role := &domain.Role{ID: rid}
		if err := db.Model(user).Association("Roles").Append(role); err != nil {
			return err
		}
	}

	return nil
}
