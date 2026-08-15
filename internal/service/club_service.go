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
	clubRepo        *repository.ClubRepository
	joinRequestRepo *repository.JoinRequestRepository
	hierarchyRepo   *repository.HierarchyRepository
}

func NewClubService(clubRepo *repository.ClubRepository, joinRequestRepo *repository.JoinRequestRepository, hierarchyRepo *repository.HierarchyRepository) *ClubService {
	return &ClubService{
		clubRepo:        clubRepo,
		joinRequestRepo: joinRequestRepo,
		hierarchyRepo:   hierarchyRepo,
	}
}

func (s *ClubService) CreateClub(name, description, joinCode string, ownerID uuid.UUID) (*model.Club, error) {
	// Check if joinCode is unique
	existingClub, _ := s.clubRepo.FindByJoinCode(joinCode)
	if existingClub != nil {
		return nil, errors.New("join code already exists, please choose another one")
	}

	club := &model.Club{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		JoinCode:    joinCode,
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

func (s *ClubService) JoinClub(joinCode string, userID uuid.UUID, domainID *uuid.UUID) error {
	club, err := s.clubRepo.FindByJoinCode(joinCode)
	if err != nil {
		return errors.New("invalid join code")
	}

	isMember, err := s.clubRepo.IsMember(club.ID, userID)
	if err != nil {
		return err
	}
	if isMember {
		return ErrAlreadyMember
	}

	// Check if already pending
	existingReq, err := s.joinRequestRepo.FindPendingByUserAndClub(userID, club.ID)
	if err == nil && existingReq != nil {
		return errors.New("you already have a pending join request for this club")
	}

	req := &model.JoinRequest{
		ID:       uuid.New(),
		ClubID:   club.ID,
		UserID:   userID,
		Status:   model.JoinRequestStatusPending,
		DomainID: domainID,
	}

	return s.joinRequestRepo.Create(req)
}

func (s *ClubService) GetPendingJoinRequests(clubID, userID uuid.UUID) ([]model.JoinRequest, error) {
	// Verify user is owner/admin
	members, err := s.clubRepo.FindMembers(clubID)
	if err != nil {
		return nil, err
	}
	isAdmin := false
	for _, m := range members {
		if m.UserID == userID && (m.Role == model.RoleOwner || m.Role == model.RoleAdmin) {
			isAdmin = true
			break
		}
	}
	if !isAdmin {
		return nil, ErrNotAuthorized
	}

	return s.joinRequestRepo.FindPendingByClubID(clubID)
}

func (s *ClubService) ApproveJoinRequest(requestID, adminUserID uuid.UUID) error {
	req, err := s.joinRequestRepo.FindByID(requestID)
	if err != nil {
		return errors.New("join request not found")
	}

	// Verify admin
	members, err := s.clubRepo.FindMembers(req.ClubID)
	if err != nil {
		return err
	}
	isAdmin := false
	for _, m := range members {
		if m.UserID == adminUserID && (m.Role == model.RoleOwner || m.Role == model.RoleAdmin) {
			isAdmin = true
			break
		}
	}
	if !isAdmin {
		return ErrNotAuthorized
	}

	if req.Status != model.JoinRequestStatusPending {
		return errors.New("request is not pending")
	}

	// Approve
	err = s.joinRequestRepo.UpdateStatus(requestID, model.JoinRequestStatusApproved)
	if err != nil {
		return err
	}

	// Find lowest hierarchy level
	levels, err := s.hierarchyRepo.FindByClubID(req.ClubID)
	var hierarchyLevelID *uuid.UUID
	if err == nil && len(levels) > 0 {
		lowestLevel := levels[len(levels)-1] // Since it's ordered by position ASC, last is lowest
		hierarchyLevelID = &lowestLevel.ID
	}

	// Add member
	member := &model.ClubMember{
		ID:               uuid.New(),
		ClubID:           req.ClubID,
		UserID:           req.UserID,
		Role:             model.RoleMember,
		DomainID:         req.DomainID,
		HierarchyLevelID: hierarchyLevelID,
	}
	return s.clubRepo.AddMember(member)
}

func (s *ClubService) RejectJoinRequest(requestID, adminUserID uuid.UUID) error {
	req, err := s.joinRequestRepo.FindByID(requestID)
	if err != nil {
		return errors.New("join request not found")
	}

	// Verify admin
	members, err := s.clubRepo.FindMembers(req.ClubID)
	if err != nil {
		return err
	}
	isAdmin := false
	for _, m := range members {
		if m.UserID == adminUserID && (m.Role == model.RoleOwner || m.Role == model.RoleAdmin) {
			isAdmin = true
			break
		}
	}
	if !isAdmin {
		return ErrNotAuthorized
	}

	if req.Status != model.JoinRequestStatusPending {
		return errors.New("request is not pending")
	}

	return s.joinRequestRepo.UpdateStatus(requestID, model.JoinRequestStatusRejected)
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
func (s *ClubService) UpdateSettings(clubID, requesterID uuid.UUID, ownerLabel, adminLabel, leadershipLabel string) error {
	members, err := s.clubRepo.FindMembers(clubID)
	if err != nil {
		return err
	}

	isOwner := false
	for _, m := range members {
		if m.UserID == requesterID && m.Role == model.RoleOwner {
			isOwner = true
			break
		}
	}

	if !isOwner {
		return ErrNotAuthorized
	}

	if ownerLabel == "" {
		ownerLabel = "Owner"
	}
	if adminLabel == "" {
		adminLabel = "Admin"
	}
	if leadershipLabel == "" {
		leadershipLabel = "Leadership"
	}

	return s.clubRepo.UpdateSettings(clubID, ownerLabel, adminLabel, leadershipLabel)
}
