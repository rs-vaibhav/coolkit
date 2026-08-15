package model

import (
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ClubID    uuid.UUID `gorm:"type:uuid;not null;index" json:"club_id"`
	AuthorID  uuid.UUID `gorm:"type:uuid;not null;index" json:"author_id"`
	Title     string    `gorm:"not null" json:"title"`
	Content   string    `gorm:"not null" json:"content"`
	Priority  string    `gorm:"default:'normal'" json:"priority"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Club   Club `gorm:"foreignKey:ClubID;constraint:OnDelete:CASCADE" json:"club,omitempty"`
	Author User `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE" json:"author,omitempty"`
}
