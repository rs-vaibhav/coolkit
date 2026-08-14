package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Club struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	OwnerID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"owner_id"`
	Owner       User           `gorm:"foreignKey:OwnerID" json:"owner,omitempty"`
	Members     []ClubMember   `gorm:"foreignKey:ClubID" json:"members,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Club) TableName() string {
	return "clubs"
}
