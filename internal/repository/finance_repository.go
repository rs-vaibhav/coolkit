package repository

import (
	"time"

	"github.com/coolkit-org/coolkit/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FinanceRepository struct {
	db *gorm.DB
}

func NewFinanceRepository(db *gorm.DB) *FinanceRepository {
	return &FinanceRepository{db: db}
}

func (r *FinanceRepository) Create(entry *model.FinanceEntry) error {
	return r.db.Create(entry).Error
}

func (r *FinanceRepository) FindByEventID(eventID uuid.UUID) ([]model.FinanceEntry, error) {
	var entries []model.FinanceEntry
	err := r.db.Where("event_id = ?", eventID).Preload("CreatedBy").Order("date DESC").Find(&entries).Error
	return entries, err
}

func (r *FinanceRepository) FindByID(id uuid.UUID) (*model.FinanceEntry, error) {
	var entry model.FinanceEntry
	err := r.db.First(&entry, id).Error
	if err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *FinanceRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.FinanceEntry{}, id).Error
}

func (r *FinanceRepository) FindApprovedByEventID(eventID uuid.UUID) ([]model.FinanceEntry, error) {
	var entries []model.FinanceEntry
	err := r.db.Where("event_id = ? AND status = ?", eventID, model.FinanceStatusApproved).Preload("CreatedBy").Order("date DESC").Find(&entries).Error
	return entries, err
}

func (r *FinanceRepository) UpdateStatus(id uuid.UUID, status string, approvedByID *uuid.UUID, approvedAt *time.Time, rejectionReason string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if approvedByID != nil {
		updates["approved_by_id"] = approvedByID
	}
	if approvedAt != nil {
		updates["approved_at"] = approvedAt
	}
	if status == model.FinanceStatusRejected && rejectionReason != "" {
		updates["rejection_reason"] = rejectionReason
	}
	return r.db.Model(&model.FinanceEntry{}).Where("id = ?", id).Updates(updates).Error
}

func (r *FinanceRepository) GetSummary(eventID uuid.UUID) (float64, float64, error) {
	var results []struct {
		Type  string
		Total float64
	}
	
	err := r.db.Model(&model.FinanceEntry{}).
		Select("type, sum(amount) as total").
		Where("event_id = ? AND status = ?", eventID, model.FinanceStatusApproved).
		Group("type").
		Scan(&results).Error

	if err != nil {
		return 0, 0, err
	}

	var totalIncome, totalExpense float64
	for _, res := range results {
		if res.Type == model.FinanceTypeIncome {
			totalIncome = res.Total
		} else if res.Type == model.FinanceTypeExpense {
			totalExpense = res.Total
		}
	}

	return totalIncome, totalExpense, nil
}
