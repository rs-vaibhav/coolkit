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
	ErrNotAuthorized = errors.New("you are not authorized to perform this action")
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

func (s *ClubService) UpdateMemberRole(clubID, requesterID, targetUserID uuid.UUID, newRole string) error {
	// Validate role
	validRoles := map[string]bool{model.RoleAdmin: true, model.RoleCoordinator: true, model.RoleMember: true}
	if !validRoles[newRole] {
		return errors.New("invalid role: must be admin, coordinator, or member")
	}

	// Can't change owner role
	if newRole == model.RoleOwner {
		return errors.New("cannot assign owner role")
	}

	// Check requester is owner or admin
	members, err := s.clubRepo.FindMembers(clubID)
	if err != nil {
		return err
	}

	requesterRole := ""
	targetRole := ""
	for _, m := range members {
		if m.UserID == requesterID {
			requesterRole = m.Role
		}
		if m.UserID == targetUserID {
			targetRole = m.Role
		}
	}

	if requesterRole != model.RoleOwner && requesterRole != model.RoleAdmin {
		return ErrNotAuthorized
	}

	// Admins can't change other admins or owner
	if requesterRole == model.RoleAdmin && (targetRole == model.RoleAdmin || targetRole == model.RoleOwner) {
		return ErrNotAuthorized
	}

	// Can't change the owner's role
	if targetRole == model.RoleOwner {
		return errors.New("cannot change the owner's role")
	}

	return s.clubRepo.UpdateMemberRole(clubID, targetUserID, newRole)
}

func (s *ClubService) RemoveMember(clubID, requesterID, targetUserID uuid.UUID) error {
	members, err := s.clubRepo.FindMembers(clubID)
	if err != nil {
		return err
	}

	requesterRole := ""
	targetRole := ""
	for _, m := range members {
		if m.UserID == requesterID {
			requesterRole = m.Role
		}
		if m.UserID == targetUserID {
			targetRole = m.Role
		}
	}

	if requesterRole != model.RoleOwner && requesterRole != model.RoleAdmin {
		return ErrNotAuthorized
	}

	if targetRole == model.RoleOwner {
		return errors.New("cannot remove the club owner")
	}

	if requesterRole == model.RoleAdmin && targetRole == model.RoleAdmin {
		return errors.New("admins cannot remove other admins")
	}

	return s.clubRepo.DeleteMember(clubID, targetUserID)
}

func (s *ClubService) LeaveClub(clubID, userID uuid.UUID) error {
	members, err := s.clubRepo.FindMembers(clubID)
	if err != nil {
		return err
	}

	for _, m := range members {
		if m.UserID == userID && m.Role == model.RoleOwner {
			return errors.New("owner cannot leave the club; transfer ownership first")
		}
	}

	return s.clubRepo.DeleteMember(clubID, userID)
}
