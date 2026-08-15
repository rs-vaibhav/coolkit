package repository

import (
	"github.com/coolkit-org/coolkit/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HierarchyRepository struct {
	db *gorm.DB
}

func NewHierarchyRepository(db *gorm.DB) *HierarchyRepository {
	return &HierarchyRepository{db: db}
}

func (r *HierarchyRepository) FindByClubID(clubID uuid.UUID) ([]model.HierarchyLevel, error) {
	var levels []model.HierarchyLevel
	err := r.db.Where("club_id = ?", clubID).Order("position ASC").Find(&levels).Error
	return levels, err
}

func (r *HierarchyRepository) BulkReplace(clubID uuid.UUID, levels []model.HierarchyLevel) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("club_id = ?", clubID).Delete(&model.HierarchyLevel{}).Error; err != nil {
			return err
		}
		if len(levels) > 0 {
			if err := tx.Create(&levels).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *HierarchyRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.HierarchyLevel{}, id).Error
}
