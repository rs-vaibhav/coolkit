package model

import (
	"time"

	"github.com/google/uuid"
)

type Resource struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ClubID      uuid.UUID `gorm:"type:uuid;not null;index" json:"club_id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Quantity    int       `gorm:"type:integer;not null;default:1" json:"quantity"`
	CreatedByID uuid.UUID `gorm:"type:uuid;not null;index" json:"created_by_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Club      *Club `gorm:"foreignKey:ClubID;constraint:OnDelete:CASCADE;" json:"club,omitempty"`
	CreatedBy *User `gorm:"foreignKey:CreatedByID;constraint:OnDelete:CASCADE;" json:"created_by,omitempty"`
}

func (Resource) TableName() string {
	return "resources"
}
