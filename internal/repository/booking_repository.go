package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/coolkit-org/coolkit/internal/model"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) Create(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r *BookingRepository) FindByID(id uuid.UUID) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.Where("id = ?", id).Preload("BookedBy").Preload("Resource").First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepository) FindByClubID(clubID uuid.UUID) ([]model.Booking, error) {
	var bookings []model.Booking
	err := r.db.Where("club_id = ?", clubID).
		Order("start_time ASC").
		Preload("BookedBy").
		Preload("Resource").
		Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) FindOverlappingBookings(resourceID uuid.UUID, startTime, endTime time.Time) ([]model.Booking, error) {
	var bookings []model.Booking
	// Overlap condition: start < booking.end_time AND end > booking.start_time
	err := r.db.Where("resource_id = ? AND status = ? AND start_time < ? AND end_time > ?",
		resourceID, model.BookingStatusApproved, endTime, startTime).
		Find(&bookings).Error
	return bookings, err
}

func (r *BookingRepository) UpdateStatus(id uuid.UUID, status string) error {
	return r.db.Model(&model.Booking{}).Where("id = ?", id).Update("status", status).Error
}

func (r *BookingRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Booking{}, id).Error
}
