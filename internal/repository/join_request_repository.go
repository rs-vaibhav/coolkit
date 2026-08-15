package repository

import (
	"github.com/coolkit-org/coolkit/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JoinRequestRepository struct {
	db *gorm.DB
}

func NewJoinRequestRepository(db *gorm.DB) *JoinRequestRepository {
	return &JoinRequestRepository{db: db}
}

func (r *JoinRequestRepository) Create(req *model.JoinRequest) error {
	return r.db.Create(req).Error
}

func (r *JoinRequestRepository) FindPendingByClubID(clubID uuid.UUID) ([]model.JoinRequest, error) {
	var requests []model.JoinRequest
	err := r.db.Where("club_id = ? AND status = ?", clubID, model.JoinRequestStatusPending).
		Preload("User").
		Order("created_at DESC").
		Find(&requests).Error
	return requests, err
}

func (r *JoinRequestRepository) FindByID(id uuid.UUID) (*model.JoinRequest, error) {
	var request model.JoinRequest
	err := r.db.Preload("User").First(&request, id).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *JoinRequestRepository) FindPendingByUserAndClub(userID, clubID uuid.UUID) (*model.JoinRequest, error) {
	var request model.JoinRequest
	err := r.db.Where("user_id = ? AND club_id = ? AND status = ?", userID, clubID, model.JoinRequestStatusPending).First(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}

func (r *JoinRequestRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&model.JoinRequest{}).Where("id = ?", id).Update("status", status).Error
}
