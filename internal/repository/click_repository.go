package repository

import (
	"fmt"

	"github.com/axellelanca/urlshortener/internal/models"
	"gorm.io/gorm"
)


type ClickRepository interface {
	CreateClick(click *models.Click) error
	GetClicksByLinkID(linkID uint) ([]models.Click, error)
	CountClicksByLinkID(linkID uint) (int, error)
}

type GormClickRepository struct {
	db *gorm.DB
}

func NewClickRepository(db *gorm.DB) *GormClickRepository {
	return &GormClickRepository{db: db}
}

func (r *GormClickRepository) CreateClick(click *models.Click) error {
	if err := r.db.Create(click).Error; err != nil {
		return fmt.Errorf("failed to create click: %w", err)
	}
	return nil
}

func (r *GormClickRepository) GetClicksByLinkID(linkID uint) ([]models.Click, error) {
	var clicks []models.Click
	if err := r.db.Where("link_id = ?", linkID).Find(&clicks).Error; err != nil {
		return nil, fmt.Errorf("failed to get clicks: %w", err)
	}
	return clicks, nil
}

func (r *GormClickRepository) CountClicksByLinkID(linkID uint) (int, error) {
	var count int64
	if err := r.db.Model(&models.Click{}).Where("link_id = ?", linkID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count clicks: %w", err)
	}
	return int(count), nil
}
