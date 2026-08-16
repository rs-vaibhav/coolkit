package repository

import (
	"github.com/coolkit-org/coolkit/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ClubRepository struct {
	db *gorm.DB
}

func NewClubRepository(db *gorm.DB) *ClubRepository {
	return &ClubRepository{db: db}
}

func (r *ClubRepository) Create(club *model.Club) error {
	return r.db.Create(club).Error
}

func (r *ClubRepository) FindByID(id uuid.UUID) (*model.Club, error) {
	var club model.Club
	err := r.db.Preload("Owner").First(&club, id).Error
	if err != nil {
		return nil, err
	}
	return &club, nil
}

func (r *ClubRepository) FindByJoinCode(code string) (*model.Club, error) {
	var club model.Club
	err := r.db.Preload("Owner").Where("join_code = ?", code).First(&club).Error
	if err != nil {
		return nil, err
	}
	return &club, nil
}

func (r *ClubRepository) FindByUserID(userID uuid.UUID) ([]model.Club, error) {
	var clubs []model.Club
	// Join through club_members to find all clubs a user belongs to
	err := r.db.Joins("JOIN club_members ON club_members.club_id = clubs.id").
		Where("club_members.user_id = ?", userID).
		Preload("Owner").
		Find(&clubs).Error
	return clubs, err
}

func (r *ClubRepository) AddMember(member *model.ClubMember) error {
	return r.db.Create(member).Error
}

func (r *ClubRepository) FindMembers(clubID uuid.UUID) ([]model.ClubMember, error) {
	var members []model.ClubMember
	err := r.db.Where("club_id = ?", clubID).Preload("User").Preload("Domain").Preload("HierarchyLevel").Find(&members).Error
	return members, err
}

func (r *ClubRepository) IsMember(clubID, userID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&model.ClubMember{}).
		Where("club_id = ? AND user_id = ?", clubID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ClubRepository) UpdateMemberRole(clubID, userID uuid.UUID, role string) error {
	return r.db.Model(&model.ClubMember{}).
		Where("club_id = ? AND user_id = ?", clubID, userID).
		Update("role", role).Error
}

func (r *ClubRepository) UpdateMemberOrganization(clubID, userID uuid.UUID, domainID, hierarchyLevelID *uuid.UUID) error {
	return r.db.Model(&model.ClubMember{}).
		Where("club_id = ? AND user_id = ?", clubID, userID).
		Updates(map[string]interface{}{
			"domain_id":          domainID,
			"hierarchy_level_id": hierarchyLevelID,
		}).Error
}

func (r *ClubRepository) DeleteMember(clubID, userID uuid.UUID) error {
	return r.db.Where("club_id = ? AND user_id = ?", clubID, userID).
		Delete(&model.ClubMember{}).Error
}
func (r *ClubRepository) UpdateSettings(clubID uuid.UUID, ownerLabel, adminLabel, leadershipLabel string) error {
	return r.db.Model(&model.Club{}).Where("id = ?", clubID).Updates(map[string]interface{}{
		"owner_label":      ownerLabel,
		"admin_label":      adminLabel,
		"leadership_label": leadershipLabel,
	}).Error
}

func (r *ClubRepository) UpdateClubImages(clubID uuid.UUID, profileImage, bannerImage *string) error {
	updates := make(map[string]interface{})
	if profileImage != nil {
		updates["profile_image"] = *profileImage
	}
	if bannerImage != nil {
		updates["banner_image"] = *bannerImage
	}
	if len(updates) == 0 {
		return nil
	}
	return r.db.Model(&model.Club{}).Where("id = ?", clubID).Updates(updates).Error
}
