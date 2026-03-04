package repository

import (
	"github.com/kenshin579/tutorials-go/database/gorm-mysql/internal/domain"
	"gorm.io/gorm"
)

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) domain.PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *domain.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) FindByID(id uint) (*domain.Post, error) {
	var post domain.Post
	err := r.db.First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) FindByUserID(userID uint) ([]domain.Post, error) {
	var posts []domain.Post
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&posts).Error
	return posts, err
}

func (r *postRepository) FindWithTags(id uint) (*domain.Post, error) {
	var post domain.Post
	err := r.db.Preload("Tags").Preload("User").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) Update(post *domain.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Post{}, id).Error
}

func (r *postRepository) AddTag(postID uint, tag *domain.Tag) error {
	var post domain.Post
	if err := r.db.First(&post, postID).Error; err != nil {
		return err
	}
	return r.db.Model(&post).Association("Tags").Append(tag)
}

func (r *postRepository) RemoveTags(postID uint) error {
	var post domain.Post
	if err := r.db.First(&post, postID).Error; err != nil {
		return err
	}
	return r.db.Model(&post).Association("Tags").Clear()
}
