package gorm_mysql

import (
	"testing"

	"github.com/kenshin579/tutorials-go/database/gorm-mysql/config"
	"github.com/kenshin579/tutorials-go/database/gorm-mysql/internal/domain"
	"github.com/kenshin579/tutorials-go/database/gorm-mysql/internal/infrastructure/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupDB(t *testing.T) *gorm.DB {
	t.Helper()
	cfg, err := config.ParseFromFile("config/config.yaml")
	require.NoError(t, err)

	db, err := database.NewMySQLDB(cfg)
	require.NoError(t, err)

	err = database.AutoMigrate(db)
	require.NoError(t, err)

	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	db.Exec("TRUNCATE TABLE post_tags")
	db.Exec("TRUNCATE TABLE tags")
	db.Exec("TRUNCATE TABLE posts")
	db.Exec("TRUNCATE TABLE profiles")
	db.Exec("TRUNCATE TABLE users")
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	return db
}

// === Phase 6: Raw SQL ===

func Test_RawSQL_Select(t *testing.T) {
	db := setupDB(t)

	db.Create(&domain.User{Name: "Frank", Email: "frank@example.com"})
	db.Create(&domain.User{Name: "Alice", Email: "alice@example.com"})

	var count int64
	err := db.Raw("SELECT COUNT(*) FROM users WHERE email LIKE ?", "%@example.com").Scan(&count).Error

	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)
}

func Test_RawSQL_Exec(t *testing.T) {
	db := setupDB(t)

	// Exec으로 벌크 INSERT
	err := db.Exec("INSERT INTO users (name, email, created_at, updated_at) VALUES (?, ?, NOW(), NOW()), (?, ?, NOW(), NOW())",
		"User1", "user1@example.com",
		"User2", "user2@example.com",
	).Error

	assert.NoError(t, err)

	var count int64
	db.Model(&domain.User{}).Count(&count)
	assert.Equal(t, int64(2), count)
}

// === Phase 6: Scopes ===

func Paginate(offset, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(limit)
	}
}

func ByName(name string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("name = ?", name)
	}
}

func Test_Scopes(t *testing.T) {
	db := setupDB(t)

	db.Create(&domain.User{Name: "Frank", Email: "frank@example.com"})
	db.Create(&domain.User{Name: "Alice", Email: "alice@example.com"})
	db.Create(&domain.User{Name: "Bob", Email: "bob@example.com"})

	// Scope 조합 사용
	var users []domain.User
	err := db.Scopes(Paginate(0, 2)).Find(&users).Error
	assert.NoError(t, err)
	assert.Len(t, users, 2)

	// 이름 필터 Scope
	var filtered []domain.User
	err = db.Scopes(ByName("Frank")).Find(&filtered).Error
	assert.NoError(t, err)
	assert.Len(t, filtered, 1)
	assert.Equal(t, "Frank", filtered[0].Name)
}

// === Phase 6: Hook ===

type AuditUser struct {
	gorm.Model
	Name  string `gorm:"type:varchar(100);not null"`
	Email string `gorm:"type:varchar(200);not null"`
	Slug  string `gorm:"type:varchar(100)"`
}

func (AuditUser) TableName() string {
	return "audit_users"
}

func (u *AuditUser) BeforeCreate(tx *gorm.DB) error {
	u.Slug = "user-" + u.Name
	return nil
}

func Test_Hook_BeforeCreate(t *testing.T) {
	db := setupDB(t)
	db.AutoMigrate(&AuditUser{})
	defer db.Exec("DROP TABLE IF EXISTS audit_users")

	user := &AuditUser{Name: "Frank", Email: "frank@example.com"}
	err := db.Create(user).Error

	assert.NoError(t, err)
	assert.Equal(t, "user-Frank", user.Slug)
}

// === Phase 7: 전체 시나리오 통합 테스트 ===

func Test_FullScenario(t *testing.T) {
	db := setupDB(t)

	// 1. User + Profile 생성 (트랜잭션)
	user := &domain.User{
		Name:  "Frank",
		Email: "frank@example.com",
		Profile: domain.Profile{
			Bio:    "Go developer",
			Avatar: "https://example.com/frank.png",
		},
	}
	err := db.Transaction(func(tx *gorm.DB) error {
		return tx.Create(user).Error
	})
	require.NoError(t, err)

	// 2. Post 생성
	post1 := &domain.Post{Title: "GORM 입문", Content: "GORM v2 사용법", UserID: user.ID}
	post2 := &domain.Post{Title: "Clean Architecture", Content: "Go 프로젝트 구조", UserID: user.ID}
	db.Create(post1)
	db.Create(post2)

	// 3. Tag 생성 및 N:M 연결
	goTag := &domain.Tag{Name: "Go"}
	dbTag := &domain.Tag{Name: "Database"}
	db.Create(goTag)
	db.Create(dbTag)

	db.Model(post1).Association("Tags").Append(goTag, dbTag)
	db.Model(post2).Association("Tags").Append(goTag)

	// 4. 전체 조회 검증
	var foundUser domain.User
	err = db.Preload("Profile").Preload("Posts.Tags").First(&foundUser, user.ID).Error
	require.NoError(t, err)

	assert.Equal(t, "Frank", foundUser.Name)
	assert.Equal(t, "Go developer", foundUser.Profile.Bio)
	assert.Len(t, foundUser.Posts, 2)

	// post1에 2개 태그, post2에 1개 태그
	for _, p := range foundUser.Posts {
		if p.Title == "GORM 입문" {
			assert.Len(t, p.Tags, 2)
		} else {
			assert.Len(t, p.Tags, 1)
		}
	}

	// 5. Soft Delete 후 조회
	db.Delete(&domain.Post{}, post2.ID)

	var activePosts []domain.Post
	db.Where("user_id = ?", user.ID).Find(&activePosts)
	assert.Len(t, activePosts, 1)
}
