package model

import (
	"time"

	"github.com/google/uuid"
)

// HierarchyLevel defines a rank/position within a club's organizational structure.
// Levels are club-wide and shared across all domains.
// Position 1 = highest rank (e.g., Head), Position 2 = next (e.g., Lead), etc.
type HierarchyLevel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ClubID    uuid.UUID `gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE" json:"club_id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Position  int       `gorm:"not null" json:"position"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Club Club `gorm:"foreignKey:ClubID" json:"club,omitempty"`
}

func (HierarchyLevel) TableName() string {
	return "hierarchy_levels"
}
