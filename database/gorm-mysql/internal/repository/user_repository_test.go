package repository

import (
	"testing"

	"github.com/kenshin579/tutorials-go/database/gorm-mysql/config"
	"github.com/kenshin579/tutorials-go/database/gorm-mysql/internal/domain"
	"github.com/kenshin579/tutorials-go/database/gorm-mysql/internal/infrastructure/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	cfg, err := config.ParseFromFile("../../config/config.yaml")
	require.NoError(t, err)

	db, err := database.NewMySQLDB(cfg)
	require.NoError(t, err)

	err = database.AutoMigrate(db)
	require.NoError(t, err)

	// 테스트 전 데이터 정리
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	db.Exec("TRUNCATE TABLE post_tags")
	db.Exec("TRUNCATE TABLE tags")
	db.Exec("TRUNCATE TABLE posts")
	db.Exec("TRUNCATE TABLE profiles")
	db.Exec("TRUNCATE TABLE users")
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")

	return db
}

// === Phase 4: CRUD 기본 조작 ===

func Test_UserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &domain.User{Name: "Frank", Email: "frank@example.com"}
	err := repo.Create(user)

	assert.NoError(t, err)
	assert.NotZero(t, user.ID)
}

func Test_UserRepository_CreateBatch(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	users := []domain.User{
		{Name: "Alice", Email: "alice@example.com"},
		{Name: "Bob", Email: "bob@example.com"},
		{Name: "Charlie", Email: "charlie@example.com"},
	}
	err := repo.CreateBatch(users)

	assert.NoError(t, err)
	for _, u := range users {
		assert.NotZero(t, u.ID)
	}
}

func Test_UserRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &domain.User{Name: "Frank", Email: "frank@example.com"}
	repo.Create(user)

	found, err := repo.FindByID(user.ID)

	assert.NoError(t, err)
	assert.Equal(t, "Frank", found.Name)
	assert.Equal(t, "frank@example.com", found.Email)
}

func Test_UserRepository_FindByEmail(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &domain.User{Name: "Frank", Email: "frank@example.com"}
	repo.Create(user)

	found, err := repo.FindByEmail("frank@example.com")

	assert.NoError(t, err)
	assert.Equal(t, "Frank", found.Name)
}

func Test_UserRepository_FindAll(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	repo.Create(&domain.User{Name: "Alice", Email: "alice@example.com"})
	repo.Create(&domain.User{Name: "Bob", Email: "bob@example.com"})
	repo.Create(&domain.User{Name: "Charlie", Email: "charlie@example.com"})

	// 페이지네이션: offset=0, limit=2
	users, err := repo.FindAll(0, 2)

	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func Test_UserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &domain.User{Name: "Frank", Email: "frank@example.com"}
	repo.Create(user)

	user.Name = "Frank Updated"
	err := repo.Update(user)
	assert.NoError(t, err)

	found, _ := repo.FindByID(user.ID)
	assert.Equal(t, "Frank Updated", found.Name)
}

func Test_UserRepository_SoftDelete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &domain.User{Name: "Frank", Email: "frank@example.com"}
	repo.Create(user)

	err := repo.Delete(user.ID)
	assert.NoError(t, err)

	// Soft Delete: 일반 조회 시 못 찾음
	_, err = repo.FindByID(user.ID)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)

	// Unscoped 조회 시 찾을 수 있음
	var deleted domain.User
	db.Unscoped().First(&deleted, user.ID)
	assert.NotNil(t, deleted.DeletedAt)
}

func Test_UserRepository_HardDelete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &domain.User{Name: "Frank", Email: "frank@example.com"}
	repo.Create(user)

	err := repo.HardDelete(user.ID)
	assert.NoError(t, err)

	// Hard Delete: Unscoped로도 못 찾음
	var count int64
	db.Unscoped().Model(&domain.User{}).Where("id = ?", user.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}

// === Phase 5: 관계 매핑 ===

func Test_UserRepository_HasOne_Profile(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	// 1:1 - User + Profile 생성
	user := &domain.User{
		Name:  "Frank",
		Email: "frank@example.com",
		Profile: domain.Profile{
			Bio:    "Go developer",
			Avatar: "https://example.com/avatar.png",
		},
	}
	err := repo.CreateWithProfile(user)
	assert.NoError(t, err)

	// Preload로 Profile 조회
	found, err := repo.FindByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Go developer", found.Profile.Bio)
	assert.Equal(t, user.ID, found.Profile.UserID)
}

func Test_UserRepository_HasMany_Posts(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)

	user := &domain.User{Name: "Frank", Email: "frank@example.com"}
	userRepo.Create(user)

	// 1:N - User의 Post 생성
	postRepo := NewPostRepository(db)
	postRepo.Create(&domain.Post{Title: "First Post", Content: "Hello World", UserID: user.ID})
	postRepo.Create(&domain.Post{Title: "Second Post", Content: "Go is great", UserID: user.ID})

	// Preload로 Posts 조회
	found, err := userRepo.FindByID(user.ID)
	assert.NoError(t, err)
	assert.Len(t, found.Posts, 2)
}

// === Phase 6: 트랜잭션 ===

func Test_UserRepository_CreateWithProfile_Transaction(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	user := &domain.User{
		Name:  "Frank",
		Email: "frank@example.com",
		Profile: domain.Profile{
			Bio: "Gopher",
		},
	}
	err := repo.CreateWithProfile(user)
	assert.NoError(t, err)

	found, _ := repo.FindByID(user.ID)
	assert.Equal(t, "Gopher", found.Profile.Bio)
}

func Test_Transaction_Rollback(t *testing.T) {
	db := setupTestDB(t)

	// 트랜잭션 내에서 에러 발생 시 롤백 확인
	err := db.Transaction(func(tx *gorm.DB) error {
		user := &domain.User{Name: "Will Rollback", Email: "rollback@example.com"}
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		// 의도적 에러 반환 → 롤백
		return gorm.ErrInvalidData
	})
	assert.Error(t, err)

	// 롤백되어 데이터가 없어야 함
	var count int64
	db.Model(&domain.User{}).Where("email = ?", "rollback@example.com").Count(&count)
	assert.Equal(t, int64(0), count)
}
