package repository

import (
	"testing"

	"github.com/kenshin579/tutorials-go/database/gorm-mysql/internal/domain"

	"github.com/stretchr/testify/assert"
)

// === Phase 4: Post CRUD ===

func Test_PostRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)

	user := &domain.User{Name: "Frank", Email: "frank@example.com"}
	userRepo.Create(user)

	post := &domain.Post{Title: "Test Post", Content: "Content here", UserID: user.ID}
	err := postRepo.Create(post)

	assert.NoError(t, err)
	assert.NotZero(t, post.ID)
}

func Test_PostRepository_FindByUserID(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)

	user := &domain.User{Name: "Frank", Email: "frank@example.com"}
	userRepo.Create(user)

	postRepo.Create(&domain.Post{Title: "Post 1", Content: "Content 1", UserID: user.ID})
	postRepo.Create(&domain.Post{Title: "Post 2", Content: "Content 2", UserID: user.ID})

	posts, err := postRepo.FindByUserID(user.ID)

	assert.NoError(t, err)
	assert.Len(t, posts, 2)
}

func Test_PostRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)

	user := &domain.User{Name: "Frank", Email: "frank@example.com"}
	userRepo.Create(user)

	post := &domain.Post{Title: "Original", Content: "Content", UserID: user.ID}
	postRepo.Create(post)

	post.Title = "Updated Title"
	err := postRepo.Update(post)
	assert.NoError(t, err)

	found, _ := postRepo.FindByID(post.ID)
	assert.Equal(t, "Updated Title", found.Title)
}

func Test_PostRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)

	user := &domain.User{Name: "Frank", Email: "frank@example.com"}
	userRepo.Create(user)

	post := &domain.Post{Title: "To Delete", Content: "Content", UserID: user.ID}
	postRepo.Create(post)

	err := postRepo.Delete(post.ID)
	assert.NoError(t, err)
}

// === Phase 5: N:M 관계 매핑 ===

func Test_PostRepository_ManyToMany_Tags(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)

	user := &domain.User{Name: "Frank", Email: "frank@example.com"}
	userRepo.Create(user)

	// Tag 생성
	goTag := &domain.Tag{Name: "Go"}
	mysqlTag := &domain.Tag{Name: "MySQL"}
	db.Create(goTag)
	db.Create(mysqlTag)

	// Post 생성
	post := &domain.Post{Title: "GORM Tutorial", Content: "Learn GORM", UserID: user.ID}
	postRepo.Create(post)

	// N:M - Post에 Tag 추가
	err := postRepo.AddTag(post.ID, goTag)
	assert.NoError(t, err)
	err = postRepo.AddTag(post.ID, mysqlTag)
	assert.NoError(t, err)

	// Preload로 Tags 조회
	found, err := postRepo.FindWithTags(post.ID)
	assert.NoError(t, err)
	assert.Len(t, found.Tags, 2)
	assert.Equal(t, "Frank", found.User.Name)
}

func Test_PostRepository_RemoveTags(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	postRepo := NewPostRepository(db)

	user := &domain.User{Name: "Frank", Email: "frank@example.com"}
	userRepo.Create(user)

	tag := &domain.Tag{Name: "Go"}
	db.Create(tag)

	post := &domain.Post{Title: "Post", Content: "Content", UserID: user.ID}
	postRepo.Create(post)
	postRepo.AddTag(post.ID, tag)

	// Tag 전체 제거
	err := postRepo.RemoveTags(post.ID)
	assert.NoError(t, err)

	found, _ := postRepo.FindWithTags(post.ID)
	assert.Len(t, found.Tags, 0)
}
