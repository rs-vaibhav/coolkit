package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	BookingStatusPending  = "pending"
	BookingStatusApproved = "approved"
	BookingStatusRejected = "rejected"
)

type Booking struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ResourceID uuid.UUID `gorm:"type:uuid;not null;index" json:"resource_id"`
	ClubID     uuid.UUID `gorm:"type:uuid;not null;index" json:"club_id"`
	BookedByID uuid.UUID `gorm:"type:uuid;not null;index" json:"booked_by_id"`
	StartTime  time.Time `gorm:"not null" json:"start_time"`
	EndTime    time.Time `gorm:"not null" json:"end_time"`
	Purpose    string    `gorm:"type:text" json:"purpose"`
	Status     string    `gorm:"type:varchar(20);not null;default:'pending'" json:"status"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Resource *Resource `gorm:"foreignKey:ResourceID;constraint:OnDelete:CASCADE;" json:"resource,omitempty"`
	Club     *Club     `gorm:"foreignKey:ClubID;constraint:OnDelete:CASCADE;" json:"club,omitempty"`
	BookedBy *User     `gorm:"foreignKey:BookedByID;constraint:OnDelete:CASCADE;" json:"booked_by,omitempty"`
}

func (Booking) TableName() string {
	return "bookings"
}
type ClubAnalytics struct {
	MemberGrowth []MonthlyMemberCount `json:"member_growth"`
	EventStats   EventAnalytics       `json:"event_stats"`
	FinanceStats FinanceAnalytics     `json:"finance_stats"`
	TaskStats    TaskAnalytics        `json:"task_stats"`
}

type MonthlyMemberCount struct {
	Month string `json:"month"` // e.g. "2026-08"
	Count int64  `json:"count"`
}

type EventAnalytics struct {
	TotalEvents     int64 `json:"total_events"`
	UpcomingEvents  int64 `json:"upcoming_events"`
	CompletedEvents int64 `json:"completed_events"`
}

type FinanceAnalytics struct {
	TotalIncome  float64        `json:"total_income"`
	TotalExpense float64        `json:"total_expense"`
	Balance      float64        `json:"balance"`
	Categories   []CategoryStat `json:"categories"`
}

type CategoryStat struct {
	Category string  `json:"category"`
	Type     string  `json:"type"` // "income" or "expense"
	Amount   float64 `json:"amount"`
}

type TaskAnalytics struct {
	Total      int64 `json:"total"`
	Todo       int64 `json:"todo"`
	InProgress int64 `json:"in_progress"`
	Done       int64 `json:"done"`
}
