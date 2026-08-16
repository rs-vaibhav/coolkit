package repository

import (
"context"
"github.com/google/uuid"
"github.com/coolkit-org/coolkit/internal/model"
"gorm.io/gorm"
)

type ClubMemberRepository interface {
Create(ctx context.Context, member *model.ClubMember) error
GetByClubAndUser(ctx context.Context, clubID, userID uuid.UUID) (*model.ClubMember, error)
GetByClubID(ctx context.Context, clubID uuid.UUID) ([]model.ClubMember, error)
Update(ctx context.Context, member *model.ClubMember) error
Delete(ctx context.Context, id uuid.UUID) error
}

type clubMemberRepository struct {
db *gorm.DB
}

func NewClubMemberRepository(db *gorm.DB) ClubMemberRepository {
return &clubMemberRepository{db: db}
}

func (r *clubMemberRepository) Create(ctx context.Context, member *model.ClubMember) error {
return r.db.WithContext(ctx).Create(member).Error
}

func (r *clubMemberRepository) GetByClubAndUser(ctx context.Context, clubID, userID uuid.UUID) (*model.ClubMember, error) {
var member model.ClubMember
err := r.db.WithContext(ctx).Where("club_id = ? AND user_id = ?", clubID, userID).First(&member).Error
if err != nil {
return nil, err
}
return &member, nil
}

func (r *clubMemberRepository) GetByClubID(ctx context.Context, clubID uuid.UUID) ([]model.ClubMember, error) {
var members []model.ClubMember
err := r.db.WithContext(ctx).Where("club_id = ?", clubID).Preload("User").Find(&members).Error
return members, err
}

func (r *clubMemberRepository) Update(ctx context.Context, member *model.ClubMember) error {
return r.db.WithContext(ctx).Save(member).Error
}

func (r *clubMemberRepository) Delete(ctx context.Context, id uuid.UUID) error {
return r.db.WithContext(ctx).Delete(&model.ClubMember{}, "id = ?", id).Error
}
