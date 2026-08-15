package repository

import (
	"github.com/coolkit-org/coolkit/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnnouncementRepository struct {
	db *gorm.DB
}

func NewAnnouncementRepository(db *gorm.DB) *AnnouncementRepository {
	return &AnnouncementRepository{db: db}
}

func (r *AnnouncementRepository) Create(announcement *model.Announcement) error {
	return r.db.Create(announcement).Error
}

func (r *AnnouncementRepository) FindByClubID(clubID uuid.UUID) ([]model.Announcement, error) {
	var announcements []model.Announcement
	err := r.db.Where("club_id = ?", clubID).Order("created_at DESC").Preload("Author").Find(&announcements).Error
	return announcements, err
}

func (r *AnnouncementRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Announcement{}, id).Error
}

func (r *AnnouncementRepository) FindByID(id uuid.UUID) (*model.Announcement, error) {
	var announcement model.Announcement
	err := r.db.First(&announcement, id).Error
	if err != nil {
		return nil, err
	}
	return &announcement, nil
}
