package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/coolkit-org/coolkit/internal/model"
	"github.com/coolkit-org/coolkit/internal/repository"
)

type BookingService struct {
	clubRepo     *repository.ClubRepository
	resourceRepo *repository.ResourceRepository
	bookingRepo  *repository.BookingRepository
}

func NewBookingService(
	clubRepo *repository.ClubRepository,
	resourceRepo *repository.ResourceRepository,
	bookingRepo *repository.BookingRepository,
) *BookingService {
	return &BookingService{
		clubRepo:     clubRepo,
		resourceRepo: resourceRepo,
		bookingRepo:  bookingRepo,
	}
}

func (s *BookingService) isAdminOrOwner(clubID, userID uuid.UUID) (bool, error) {
	members, err := s.clubRepo.FindMembers(clubID)
	if err != nil {
		return false, err
	}
	for _, m := range members {
		if m.UserID == userID && (m.Role == model.RoleOwner || m.Role == model.RoleAdmin) {
			return true, nil
		}
	}
	return false, nil
}

func (s *BookingService) AddResource(name, description string, quantity int, clubID, userID uuid.UUID) (*model.Resource, error) {
	isAdmin, err := s.isAdminOrOwner(clubID, userID)
	if err != nil {
		return nil, err
	}
	if !isAdmin {
		return nil, ErrNotAuthorized
	}

	if quantity < 1 {
		return nil, errors.New("quantity must be at least 1")
	}

	resource := &model.Resource{
		ID:          uuid.New(),
		ClubID:      clubID,
		Name:        name,
		Description: description,
		Quantity:    quantity,
		CreatedByID: userID,
	}

	err = s.resourceRepo.Create(resource)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

func (s *BookingService) ListResources(clubID uuid.UUID) ([]model.Resource, error) {
	return s.resourceRepo.FindByClubID(clubID)
}

func (s *BookingService) DeleteResource(resourceID, userID uuid.UUID) error {
	resource, err := s.resourceRepo.FindByID(resourceID)
	if err != nil {
		return errors.New("resource not found")
	}

	isAdmin, err := s.isAdminOrOwner(resource.ClubID, userID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return ErrNotAuthorized
	}

	return s.resourceRepo.Delete(resourceID)
}

func (s *BookingService) BookResource(resourceID, clubID, userID uuid.UUID, startTime, endTime time.Time, purpose string) (*model.Booking, error) {
	// Verify membership
	isMem, err := s.clubRepo.IsMember(clubID, userID)
	if err != nil {
		return nil, err
	}
	if !isMem {
		return nil, ErrNotAuthorized
	}

	// Validate inputs
	if !startTime.Before(endTime) {
		return nil, errors.New("start time must be before end time")
	}
	if startTime.Before(time.Now()) {
		return nil, errors.New("cannot book in the past")
	}

	_, err = s.resourceRepo.FindByID(resourceID)
	if err != nil {
		return nil, errors.New("resource not found")
	}

	booking := &model.Booking{
		ID:         uuid.New(),
		ResourceID: resourceID,
		ClubID:     clubID,
		BookedByID: userID,
		StartTime:  startTime,
		EndTime:    endTime,
		Purpose:    purpose,
		Status:     model.BookingStatusPending,
	}

	err = s.bookingRepo.Create(booking)
	if err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *BookingService) ListClubBookings(clubID uuid.UUID) ([]model.Booking, error) {
	return s.bookingRepo.FindByClubID(clubID)
}

func (s *BookingService) ApproveBooking(bookingID, userID uuid.UUID) error {
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return errors.New("booking not found")
	}

	isAdmin, err := s.isAdminOrOwner(booking.ClubID, userID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return ErrNotAuthorized
	}

	if booking.Status != model.BookingStatusPending {
		return errors.New("booking is not pending")
	}

	// Check for overlap limits
	resource, err := s.resourceRepo.FindByID(booking.ResourceID)
	if err != nil {
		return errors.New("associated resource not found")
	}

	overlapping, err := s.bookingRepo.FindOverlappingBookings(booking.ResourceID, booking.StartTime, booking.EndTime)
	if err != nil {
		return err
	}

	// Calculate concurrently approved bookings
	// Since we only allow exact overlaps, let's count overlapping instances.
	// For maximum robustness, we find the maximum overlap at any single point in time.
	// Let's implement an interval coverage count.
	maxOverlap := 0
	for _, o := range overlapping {
		if o.ID == booking.ID {
			continue
		}
		// Count how many other approved bookings overlap with this booking's timeframe
		overlapCount := 1 // start with o itself
		for _, other := range overlapping {
			if other.ID == booking.ID || other.ID == o.ID {
				continue
			}
			// check if they overlap
			if other.StartTime.Before(o.EndTime) && other.EndTime.After(o.StartTime) {
				overlapCount++
			}
		}
		if overlapCount > maxOverlap {
			maxOverlap = overlapCount
		}
	}

	if maxOverlap >= resource.Quantity {
		return errors.New("resource is fully booked during this time slot")
	}

	return s.bookingRepo.UpdateStatus(bookingID, model.BookingStatusApproved)
}

func (s *BookingService) RejectBooking(bookingID, userID uuid.UUID) error {
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return errors.New("booking not found")
	}

	isAdmin, err := s.isAdminOrOwner(booking.ClubID, userID)
	if err != nil {
		return err
	}
	if !isAdmin {
		return ErrNotAuthorized
	}

	if booking.Status != model.BookingStatusPending {
		return errors.New("booking is not pending")
	}

	return s.bookingRepo.UpdateStatus(bookingID, model.BookingStatusRejected)
}

func (s *BookingService) CancelBooking(bookingID, userID uuid.UUID) error {
	booking, err := s.bookingRepo.FindByID(bookingID)
	if err != nil {
		return errors.New("booking not found")
	}

	// Booker or Admin can cancel/delete
	isAdmin, err := s.isAdminOrOwner(booking.ClubID, userID)
	if err != nil {
		return err
	}

	if booking.BookedByID != userID && !isAdmin {
		return ErrNotAuthorized
	}

	return s.bookingRepo.Delete(bookingID)
}
