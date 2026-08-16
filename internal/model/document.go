package model

import (
"time"

"github.com/google/uuid"
"gorm.io/gorm"
)

type DocumentType string

const (
DocumentPermissionLetter DocumentType = "permission_letter"
DocumentOther            DocumentType = "other"
)

type Document struct {
ID          uuid.UUID    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
ClubID      uuid.UUID    `gorm:"type:uuid;not null;index" json:"club_id"`
EventID     *uuid.UUID   `gorm:"type:uuid;index" json:"event_id,omitempty"`
Title       string       `gorm:"not null" json:"title"`
Description string       `json:"description"`
DocType     DocumentType `gorm:"type:varchar(50);not null" json:"doc_type"`
FileURL     string       `gorm:"type:varchar(500);not null" json:"file_url"`
UploadedBy  uuid.UUID    `gorm:"type:uuid;not null" json:"uploaded_by"`
CreatedAt   time.Time    `gorm:"autoCreateTime" json:"created_at"`
UpdatedAt   time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

Club Club `gorm:"foreignKey:ClubID" json:"club,omitempty"`
}

func (Document) TableName() string {
return "documents"
}
