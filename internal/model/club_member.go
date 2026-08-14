package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	RoleOwner       = "owner"
	RoleAdmin       = "admin"
	RoleCoordinator = "coordinator"
	RoleMember      = "member"
)

type ClubMember struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ClubID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_club_user" json:"club_id"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_club_user" json:"user_id"`
	Role     string    `gorm:"not null;default:'member'" json:"role"`
	Club     Club      `gorm:"foreignKey:ClubID" json:"club,omitempty"`
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	JoinedAt time.Time `gorm:"autoCreateTime" json:"joined_at"`
}

func (ClubMember) TableName() string {
	return "club_members"
}
