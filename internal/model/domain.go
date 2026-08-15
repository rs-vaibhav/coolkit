package model

import (
	"time"

	"github.com/google/uuid"
)

// Domain represents an organizational domain within a club (e.g., "Creative", "Development", "Marketing").
// Members are assigned to exactly one domain.
type Domain struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ClubID      uuid.UUID `gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE" json:"club_id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description string    `gorm:"type:varchar(500)" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`

	Club Club `gorm:"foreignKey:ClubID" json:"club,omitempty"`
}

func (Domain) TableName() string {
	return "domains"
}
