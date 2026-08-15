package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	JoinRequestStatusPending  = "pending"
	JoinRequestStatusApproved = "approved"
	JoinRequestStatusRejected = "rejected"
)

type JoinRequest struct {
	ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ClubID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"club_id"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null;index" json:"user_id"`
	DomainID  *uuid.UUID `gorm:"type:uuid;index" json:"domain_id"`
	Status    string     `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	Club   Club    `gorm:"foreignKey:ClubID;constraint:OnDelete:CASCADE" json:"club,omitempty"`
	User   User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Domain *Domain `gorm:"foreignKey:DomainID" json:"domain,omitempty"`
}
