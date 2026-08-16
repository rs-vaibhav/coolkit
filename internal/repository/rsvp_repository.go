package repository

import (
"context"
"github.com/google/uuid"
"github.com/coolkit-org/coolkit/internal/model"
"gorm.io/gorm"
)

type RSVPRepository interface {
Create(ctx context.Context, rsvp *model.RSVP) error
GetByID(ctx context.Context, id uuid.UUID) (*model.RSVP, error)
GetByEventAndUser(ctx context.Context, eventID, userID uuid.UUID) (*model.RSVP, error)
GetByEventID(ctx context.Context, eventID uuid.UUID) ([]model.RSVP, error)
Update(ctx context.Context, rsvp *model.RSVP) error
Delete(ctx context.Context, id uuid.UUID) error
GetCountByStatus(ctx context.Context, eventID uuid.UUID, status model.RSVPStatus) (int64, error)
}

type rsvpRepository struct {
db *gorm.DB
}

func NewRSVPRepository(db *gorm.DB) RSVPRepository {
return &rsvpRepository{db: db}
}

func (r *rsvpRepository) Create(ctx context.Context, rsvp *model.RSVP) error {
return r.db.WithContext(ctx).Create(rsvp).Error
}

func (r *rsvpRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.RSVP, error) {
var rsvp model.RSVP
err := r.db.WithContext(ctx).Preload("User").First(&rsvp, "id = ?", id).Error
if err != nil {
return nil, err
}
return &rsvp, nil
}

func (r *rsvpRepository) GetByEventAndUser(ctx context.Context, eventID, userID uuid.UUID) (*model.RSVP, error) {
var rsvp model.RSVP
err := r.db.WithContext(ctx).Where("event_id = ? AND user_id = ?", eventID, userID).First(&rsvp).Error
if err != nil {
return nil, err
}
return &rsvp, nil
}

func (r *rsvpRepository) GetByEventID(ctx context.Context, eventID uuid.UUID) ([]model.RSVP, error) {
var rsvps []model.RSVP
err := r.db.WithContext(ctx).Preload("User").Where("event_id = ?", eventID).Find(&rsvps).Error
return rsvps, err
}

func (r *rsvpRepository) Update(ctx context.Context, rsvp *model.RSVP) error {
return r.db.WithContext(ctx).Save(rsvp).Error
}

func (r *rsvpRepository) Delete(ctx context.Context, id uuid.UUID) error {
return r.db.WithContext(ctx).Delete(&model.RSVP{}, "id = ?", id).Error
}

func (r *rsvpRepository) GetCountByStatus(ctx context.Context, eventID uuid.UUID, status model.RSVPStatus) (int64, error) {
var count int64
err := r.db.WithContext(ctx).Model(&model.RSVP{}).Where("event_id = ? AND status = ?", eventID, status).Count(&count).Error
return count, err
}
