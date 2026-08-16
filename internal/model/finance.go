package model

import (
	"time"

	"github.com/google/uuid"
)

const (
	FinanceTypeIncome  = "income"
	FinanceTypeExpense = "expense"
)

const (
	FinanceStatusPending   = "pending"
	FinanceStatusApproved  = "approved"
	FinanceStatusRejected  = "rejected"
)

type FinanceEntry struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	EventID           uuid.UUID `gorm:"type:uuid;not null;index" json:"event_id"`
	CreatedByID       uuid.UUID `gorm:"type:uuid;not null;index" json:"created_by_id"`
	Type              string    `gorm:"type:varchar(20);not null" json:"type"`
	Category          string    `gorm:"type:varchar(100);not null" json:"category"`
	Amount            float64   `gorm:"not null" json:"amount"`
	Description       string    `json:"description"`
	Date              time.Time `gorm:"not null" json:"date"`
	ProofImage        string    `gorm:"type:varchar(500)" json:"proof_image"`
	Status            string    `gorm:"type:varchar(20);default:'pending'" json:"status"`
	ApprovedByID      *uuid.UUID `gorm:"type:uuid;index" json:"approved_by_id,omitempty"`
	ApprovedAt        *time.Time `json:"approved_at,omitempty"`
	RejectionReason   string    `json:"rejection_reason,omitempty"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Event     Event `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE;" json:"event"`
	CreatedBy User  `gorm:"foreignKey:CreatedByID;constraint:OnDelete:CASCADE;" json:"created_by"`
	ApprovedBy *User  `gorm:"foreignKey:ApprovedByID;constraint:OnDelete:SET NULL;" json:"approved_by,omitempty"`
}

type FinanceSummary struct {
	TotalIncome  float64        `json:"total_income"`
	TotalExpense float64        `json:"total_expense"`
	Balance      float64        `json:"balance"`
	Entries      []FinanceEntry `json:"entries"`
}
