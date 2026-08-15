package repository

import (
	"github.com/coolkit-org/coolkit/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) Create(event *model.Event) error {
	return r.db.Create(event).Error
}

func (r *EventRepository) FindByClubID(clubID uuid.UUID) ([]model.Event, error) {
	var events []model.Event
	err := r.db.Where("club_id = ?", clubID).Order("date ASC").Find(&events).Error
	return events, err
}

func (r *EventRepository) FindByID(id uuid.UUID) (*model.Event, error) {
	var event model.Event
	err := r.db.First(&event, id).Error
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *EventRepository) Update(event *model.Event) error {
	return r.db.Save(event).Error
}

func (r *EventRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Event{}, id).Error
}
