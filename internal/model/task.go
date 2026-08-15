package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	TaskStatusTodo       = "todo"
	TaskStatusInProgress = "in_progress"
	TaskStatusDone       = "done"
)

type Task struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	EventID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"event_id"`
	Event        *Event     `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE;" json:"event,omitempty"`
	AssignedToID uuid.UUID  `gorm:"type:uuid;not null;index" json:"assigned_to_id"`
	AssignedTo   *User      `gorm:"foreignKey:AssignedToID;constraint:OnDelete:CASCADE;" json:"assigned_to,omitempty"`
	Title        string     `gorm:"not null" json:"title"`
	Description  string     `json:"description"`
	Status       string     `gorm:"type:varchar(20);not null;default:'todo'" json:"status"`
	DueDate      *time.Time `json:"due_date"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
