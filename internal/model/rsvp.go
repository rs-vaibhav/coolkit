package model

import (
"time"

"github.com/google/uuid"
)

type RSVPStatus string

const (
RSVPStatusGoing       RSVPStatus = "going"
RSVPStatusInterested  RSVPStatus = "interested"
RSVPStatusNotGoing    RSVPStatus = "not_going"
)

type RSVP struct {
ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
EventID   uuid.UUID `gorm:"type:uuid;not null;index:idx_event_user" json:"event_id"`
UserID    uuid.UUID `gorm:"type:uuid;not null;index:idx_event_user" json:"user_id"`
Status    RSVPStatus `gorm:"type:varchar(20);not null" json:"status"`
CheckedIn bool      `gorm:"default:false" json:"checked_in"`
CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

Event Event `gorm:"foreignKey:EventID" json:"event,omitempty"`
User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (RSVP) TableName() string {
return "rsvps"
}
