package model

import (
	"time"

	"github.com/google/uuid"
)

type EventRole struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EventID   uuid.UUID `gorm:"type:uuid;not null;index" json:"event_id"`
	Event     Event     `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE;" json:"event,omitempty"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"user,omitempty"`
	RoleName  string    `gorm:"type:varchar(100);not null" json:"role_name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
