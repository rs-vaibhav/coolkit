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
		var existingLevels []model.HierarchyLevel
		if err := tx.Where("club_id = ?", clubID).Find(&existingLevels).Error; err != nil {
			return err
		}

		existingByPosition := make(map[int]uuid.UUID)
		for _, l := range existingLevels {
			existingByPosition[l.Position] = l.ID
		}

		var keepIDs []uuid.UUID
		for _, l := range levels {
			if id, exists := existingByPosition[l.Position]; exists {
				if err := tx.Model(&model.HierarchyLevel{}).Where("id = ?", id).Update("name", l.Name).Error; err != nil {
					return err
				}
				keepIDs = append(keepIDs, id)
			} else {
				if err := tx.Create(&l).Error; err != nil {
					return err
				}
				keepIDs = append(keepIDs, l.ID)
			}
		}

		if len(keepIDs) > 0 {
			if err := tx.Model(&model.ClubMember{}).Where("club_id = ? AND hierarchy_level_id NOT IN ?", clubID, keepIDs).Update("hierarchy_level_id", nil).Error; err != nil {
				return err
			}
			if err := tx.Where("club_id = ? AND id NOT IN ?", clubID, keepIDs).Delete(&model.HierarchyLevel{}).Error; err != nil {
				return err
			}
		} else {
			if err := tx.Model(&model.ClubMember{}).Where("club_id = ?", clubID).Update("hierarchy_level_id", nil).Error; err != nil {
				return err
			}
			if err := tx.Where("club_id = ?", clubID).Delete(&model.HierarchyLevel{}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *HierarchyRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.HierarchyLevel{}, id).Error
}
