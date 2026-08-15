package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/coolkit-org/coolkit/internal/model"
)

type ResourceRepository struct {
	db *gorm.DB
}

func NewResourceRepository(db *gorm.DB) *ResourceRepository {
	return &ResourceRepository{db: db}
}

func (r *ResourceRepository) Create(resource *model.Resource) error {
	return r.db.Create(resource).Error
}

func (r *ResourceRepository) FindByID(id uuid.UUID) (*model.Resource, error) {
	var res model.Resource
	err := r.db.Where("id = ?", id).First(&res).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *ResourceRepository) FindByClubID(clubID uuid.UUID) ([]model.Resource, error) {
	var resources []model.Resource
	err := r.db.Where("club_id = ?", clubID).Order("created_at DESC").Find(&resources).Error
	return resources, err
}

func (r *ResourceRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Resource{}, id).Error
}
