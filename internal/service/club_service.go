package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/coolkit-org/coolkit/internal/model"
	"github.com/coolkit-org/coolkit/internal/repository"
)

var (
	ErrAlreadyMember = errors.New("user is already a member of this club")
	ErrClubNotFound  = errors.New("club not found")
)

type ClubService struct {
	clubRepo *repository.ClubRepository
}

func NewClubService(clubRepo *repository.ClubRepository) *ClubService {
	return &ClubService{
		clubRepo: clubRepo,
	}
}

func (s *ClubService) CreateClub(name, description string, ownerID uuid.UUID) (*model.Club, error) {
	club := &model.Club{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
	}

	err := s.clubRepo.Create(club)
	if err != nil {
		return nil, err
	}

	member := &model.ClubMember{
		ID:     uuid.New(),
		ClubID: club.ID,
		UserID: ownerID,
		Role:   model.RoleOwner,
	}
	err = s.clubRepo.AddMember(member)
	if err != nil {
		return nil, err
	}

	return club, nil
}

func (s *ClubService) GetUserClubs(userID uuid.UUID) ([]model.Club, error) {
	return s.clubRepo.FindByUserID(userID)
}

func (s *ClubService) GetClub(clubID uuid.UUID) (*model.Club, error) {
	club, err := s.clubRepo.FindByID(clubID)
	if err != nil {
		return nil, ErrClubNotFound
	}
	return club, nil
}

func (s *ClubService) JoinClub(clubID, userID uuid.UUID) error {
	isMember, err := s.clubRepo.IsMember(clubID, userID)
	if err != nil {
		return err
	}
	if isMember {
		return ErrAlreadyMember
	}

	member := &model.ClubMember{
		ID:     uuid.New(),
		ClubID: clubID,
		UserID: userID,
		Role:   model.RoleMember,
	}
	return s.clubRepo.AddMember(member)
}

func (s *ClubService) GetClubMembers(clubID uuid.UUID) ([]model.ClubMember, error) {
	return s.clubRepo.FindMembers(clubID)
}
