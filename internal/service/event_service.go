package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/coolkit-org/coolkit/internal/model"
	"github.com/coolkit-org/coolkit/internal/pkg/formbricks"
	"github.com/coolkit-org/coolkit/internal/repository"
	"github.com/google/uuid"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrUnauthorized  = errors.New("user is not authorized to perform this action")
)

type EventService struct {
	eventRepo     *repository.EventRepository
	clubRepo      *repository.ClubRepository
	eventRoleRepo *repository.EventRoleRepository
}

func NewEventService(eventRepo *repository.EventRepository, clubRepo *repository.ClubRepository, eventRoleRepo *repository.EventRoleRepository) *EventService {
	return &EventService{
		eventRepo:     eventRepo,
		clubRepo:      clubRepo,
		eventRoleRepo: eventRoleRepo,
	}
}

func (s *EventService) CreateEvent(clubID, userID uuid.UUID, title, description string, date time.Time, location string) (*model.Event, error) {
	// Verify user is an owner or admin of the club
	members, err := s.clubRepo.FindMembers(clubID)
	if err != nil {
		return nil, err
	}

	isAuthorized := false
	for _, m := range members {
		if m.UserID == userID {
			if m.Role == model.RoleOwner || m.Role == model.RoleAdmin {
				isAuthorized = true
				break
			}
		}
	}

	if !isAuthorized {
		return nil, ErrUnauthorized
	}

	event := &model.Event{
		ID:          uuid.New(),
		ClubID:      clubID,
		Title:       title,
		Description: description,
		Date:        date,
		Location:    location,
	}

	err = s.eventRepo.Create(event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *EventService) GetClubEvents(clubID uuid.UUID) ([]model.Event, error) {
	return s.eventRepo.FindByClubID(clubID)
}

func (s *EventService) GetEvent(eventID uuid.UUID) (*model.Event, error) {
	return s.eventRepo.FindByID(eventID)
}

func (s *EventService) GetEventRoles(eventID uuid.UUID) ([]model.EventRole, error) {
	return s.eventRoleRepo.FindByEventID(eventID)
}

func (s *EventService) AssignEventRole(eventID, assignerID, assignedUserID uuid.UUID, roleName string) (*model.EventRole, error) {
	// 1. Get Event
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, err
	}

	// 2. Verify assigner is an owner or admin of the club
	members, err := s.clubRepo.FindMembers(event.ClubID)
	if err != nil {
		return nil, err
	}

	isAuthorized := false
	for _, m := range members {
		if m.UserID == assignerID {
			if m.Role == model.RoleOwner || m.Role == model.RoleAdmin {
				isAuthorized = true
				break
			}
		}
	}

	if !isAuthorized {
		return nil, ErrUnauthorized
	}

	// 3. Verify assigned user is in the club
	isMember := false
	for _, m := range members {
		if m.UserID == assignedUserID {
			isMember = true
			break
		}
	}
	if !isMember {
		return nil, errors.New("assigned user is not a member of the club")
	}

	// 4. Create the role
	role := &model.EventRole{
		ID:       uuid.New(),
		EventID:  eventID,
		UserID:   assignedUserID,
		RoleName: roleName,
	}

	err = s.eventRoleRepo.Create(role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *EventService) UpdateEvent(eventID, userID uuid.UUID, title, description string, date time.Time, location string) (*model.Event, error) {
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, ErrEventNotFound
	}

	// Verify user is owner/admin of the club
	if !s.isClubAdmin(event.ClubID, userID) {
		return nil, ErrUnauthorized
	}

	event.Title = title
	event.Description = description
	event.Date = date
	event.Location = location

	err = s.eventRepo.Update(event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *EventService) CreateFormbricksSurvey(eventID, userID uuid.UUID, surveyName, surveyDescription string, useDefaults bool) (*model.Event, error) {
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, ErrEventNotFound
	}

	// Verify user is owner/admin of the club
	if !s.isClubAdmin(event.ClubID, userID) {
		return nil, ErrUnauthorized
	}

	// Initialize Formbricks client
	fbClient := formbricks.NewClient()

	// Get the organization ID first
	orgID, err := fbClient.GetOrganization()
	if err != nil {
		return nil, fmt.Errorf("failed to get Formbricks organization: %w", err)
	}

	// Create a new environment for this event in the organization
	envID, err := fbClient.CreateEnvironment(orgID, event.Title)
	if err != nil {
		return nil, fmt.Errorf("failed to create Formbricks environment: %w", err)
	}

	// Create survey - either with defaults or custom
	var survey *formbricks.Survey
	if useDefaults {
		survey, err = fbClient.CreateSurveyWithDefaults(surveyName, envID)
	} else {
		survey, err = fbClient.CreateSurvey(surveyName, surveyDescription, envID)
	}
	
	if err != nil {
		return nil, fmt.Errorf("failed to create Formbricks survey: %w", err)
	}

	// Update event with Formbricks IDs
	event.FormbricksEnvID = envID
	event.FormbricksSurveyID = survey.ID
	
	err = s.eventRepo.Update(event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *EventService) UpdateEventFormbricks(eventID, userID uuid.UUID, formbricksEnvID, formbricksSurveyID string) (*model.Event, error) {
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return nil, ErrEventNotFound
	}

	// Verify user is owner/admin of the club
	if !s.isClubAdmin(event.ClubID, userID) {
		return nil, ErrUnauthorized
	}

	event.FormbricksEnvID = formbricksEnvID
	event.FormbricksSurveyID = formbricksSurveyID

	err = s.eventRepo.Update(event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *EventService) DeleteEvent(eventID, userID uuid.UUID) error {
	event, err := s.eventRepo.FindByID(eventID)
	if err != nil {
		return ErrEventNotFound
	}

	if !s.isClubAdmin(event.ClubID, userID) {
		return ErrUnauthorized
	}

	return s.eventRepo.Delete(eventID)
}

func (s *EventService) RemoveEventRole(roleID, userID uuid.UUID) error {
	role, err := s.eventRoleRepo.FindByID(roleID)
	if err != nil {
		return errors.New("role not found")
	}

	event, err := s.eventRepo.FindByID(role.EventID)
	if err != nil {
		return ErrEventNotFound
	}

	if !s.isClubAdmin(event.ClubID, userID) {
		return ErrUnauthorized
	}

	return s.eventRoleRepo.Delete(roleID)
}

// isClubAdmin checks if a user is an owner or admin of the given club
func (s *EventService) isClubAdmin(clubID, userID uuid.UUID) bool {
	members, err := s.clubRepo.FindMembers(clubID)
	if err != nil {
		return false
	}
	for _, m := range members {
		if m.UserID == userID && (m.Role == model.RoleOwner || m.Role == model.RoleAdmin) {
			return true
		}
	}
	return false
}
