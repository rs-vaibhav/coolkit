package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/coolkit-org/coolkit/internal/model"
)

type EventRoleRepository struct {
	db *gorm.DB
}

func NewEventRoleRepository(db *gorm.DB) *EventRoleRepository {
	return &EventRoleRepository{db: db}
}

func (r *EventRoleRepository) Create(role *model.EventRole) error {
	return r.db.Create(role).Error
}

func (r *EventRoleRepository) FindByEventID(eventID uuid.UUID) ([]model.EventRole, error) {
	var roles []model.EventRole
	err := r.db.Where("event_id = ?", eventID).Preload("User").Find(&roles).Error
	return roles, err
}

func (r *EventRoleRepository) FindByEventAndUser(eventID, userID uuid.UUID) (*model.EventRole, error) {
	var role model.EventRole
	err := r.db.Where("event_id = ? AND user_id = ?", eventID, userID).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *EventRoleRepository) FindByID(id uuid.UUID) (*model.EventRole, error) {
	var role model.EventRole
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *EventRoleRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.EventRole{}, id).Error
}
