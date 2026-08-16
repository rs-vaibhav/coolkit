package service

import (
"context"
"errors"
"github.com/google/uuid"
"github.com/coolkit-org/coolkit/internal/model"
"github.com/coolkit-org/coolkit/internal/repository"
)

var (
ErrRSVPNotFound     = errors.New("rsvp not found")
ErrRSVPAlreadyExists = errors.New("user already registered for this event")
)

type RSVPService interface {
CreateOrUpdateRSVP(ctx context.Context, rsvp *model.RSVP) error
GetRSVP(ctx context.Context, id uuid.UUID) (*model.RSVP, error)
GetUserRSVP(ctx context.Context, eventID, userID uuid.UUID) (*model.RSVP, error)
GetEventRSVPs(ctx context.Context, eventID uuid.UUID) ([]model.RSVP, error)
CheckInUser(ctx context.Context, rsvpID uuid.UUID) error
GetRSVPCounts(ctx context.Context, eventID uuid.UUID) (map[string]int64, error)
DeleteRSVP(ctx context.Context, id uuid.UUID) error
}

type rsvpService struct {
rsvpRepo   repository.RSVPRepository
memberRepo repository.ClubMemberRepository
eventRepo  *repository.EventRepository
}

func NewRSVPService(rsvpRepo repository.RSVPRepository, memberRepo repository.ClubMemberRepository, eventRepo *repository.EventRepository) RSVPService {
return &rsvpService{rsvpRepo: rsvpRepo, memberRepo: memberRepo, eventRepo: eventRepo}
}

func (s *rsvpService) CreateOrUpdateRSVP(ctx context.Context, rsvp *model.RSVP) error {
// Check if event exists
_, err := s.eventRepo.FindByID(rsvp.EventID)
if err != nil {
return errors.New("event not found")
}

// Check if user already has RSVP
existing, err := s.rsvpRepo.GetByEventAndUser(ctx, rsvp.EventID, rsvp.UserID)
if err == nil {
// Update existing
existing.Status = rsvp.Status
return s.rsvpRepo.Update(ctx, existing)
}

// Create new
return s.rsvpRepo.Create(ctx, rsvp)
}

func (s *rsvpService) GetRSVP(ctx context.Context, id uuid.UUID) (*model.RSVP, error) {
rsvp, err := s.rsvpRepo.GetByID(ctx, id)
if err != nil {
return nil, ErrRSVPNotFound
}
return rsvp, nil
}

func (s *rsvpService) GetUserRSVP(ctx context.Context, eventID, userID uuid.UUID) (*model.RSVP, error) {
rsvp, err := s.rsvpRepo.GetByEventAndUser(ctx, eventID, userID)
if err != nil {
return nil, ErrRSVPNotFound
}
return rsvp, nil
}

func (s *rsvpService) GetEventRSVPs(ctx context.Context, eventID uuid.UUID) ([]model.RSVP, error) {
return s.rsvpRepo.GetByEventID(ctx, eventID)
}

func (s *rsvpService) CheckInUser(ctx context.Context, rsvpID uuid.UUID) error {
rsvp, err := s.rsvpRepo.GetByID(ctx, rsvpID)
if err != nil {
return ErrRSVPNotFound
}

rsvp.CheckedIn = true
return s.rsvpRepo.Update(ctx, rsvp)
}

func (s *rsvpService) GetRSVPCounts(ctx context.Context, eventID uuid.UUID) (map[string]int64, error) {
counts := make(map[string]int64)

going, _ := s.rsvpRepo.GetCountByStatus(ctx, eventID, model.RSVPStatusGoing)
interested, _ := s.rsvpRepo.GetCountByStatus(ctx, eventID, model.RSVPStatusInterested)
notGoing, _ := s.rsvpRepo.GetCountByStatus(ctx, eventID, model.RSVPStatusNotGoing)

counts["going"] = going
counts["interested"] = interested
counts["not_going"] = notGoing
counts["total"] = going + interested + notGoing

return counts, nil
}

func (s *rsvpService) DeleteRSVP(ctx context.Context, id uuid.UUID) error {
_, err := s.rsvpRepo.GetByID(ctx, id)
if err != nil {
return ErrRSVPNotFound
}
return s.rsvpRepo.Delete(ctx, id)
}
