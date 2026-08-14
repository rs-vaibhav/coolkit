package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name         string         `gorm:"not null" json:"name"`
	Email        string         `gorm:"unique;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	AvatarURL    string         `json:"avatar_url"`
	OwnedClubs   []Club         `gorm:"foreignKey:OwnerID" json:"owned_clubs,omitempty"`
	Memberships  []ClubMember   `gorm:"foreignKey:UserID" json:"memberships,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}
