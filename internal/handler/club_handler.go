package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/coolkit-org/coolkit/internal/service"
	"github.com/coolkit-org/coolkit/pkg/response"
)

type ClubHandler struct {
	clubService *service.ClubService
}

func NewClubHandler(clubService *service.ClubService) *ClubHandler {
	return &ClubHandler{clubService: clubService}
}

type CreateClubRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	JoinCode    string `json:"join_code" binding:"required,min=4,max=50"`
}

func (h *ClubHandler) Create(c *gin.Context) {
	var req CreateClubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User ID not found in context")
		return
	}

	userID := userIDVal.(uuid.UUID)

	club, err := h.clubService.CreateClub(req.Name, req.Description, req.JoinCode, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, club)
}

func (h *ClubHandler) List(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User ID not found in context")
		return
	}

	userID := userIDVal.(uuid.UUID)

	clubs, err := h.clubService.GetUserClubs(userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, clubs)
}

func (h *ClubHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	clubID, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid club ID")
		return
	}

	club, err := h.clubService.GetClub(clubID)
	if err != nil {
		if errors.Is(err, service.ErrClubNotFound) {
			response.NotFound(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, club)
}

type JoinClubRequest struct {
	JoinCode string  `json:"join_code" binding:"required"`
	DomainID *string `json:"domain_id"`
}

func (h *ClubHandler) Join(c *gin.Context) {
	var req JoinClubRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User ID not found in context")
		return
	}

	userID := userIDVal.(uuid.UUID)

	var domainID *uuid.UUID
	if req.DomainID != nil && *req.DomainID != "" {
		id, err := uuid.Parse(*req.DomainID)
		if err == nil {
			domainID = &id
		}
	}

	err := h.clubService.JoinClub(req.JoinCode, userID, domainID)
	if err != nil {
		if errors.Is(err, service.ErrAlreadyMember) || err.Error() == "you already have a pending join request for this club" {
			response.Error(c, 409, err.Error())
			return
		}
		if err.Error() == "invalid join code" {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Join request sent! You will be added to the club once the owner approves it."})
}

func (h *ClubHandler) Members(c *gin.Context) {
	idStr := c.Param("id")
	clubID, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid club ID")
		return
	}

	members, err := h.clubService.GetClubMembers(clubID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, members)
}

func (h *ClubHandler) GetJoinRequests(c *gin.Context) {
	idStr := c.Param("id")
	clubID, err := uuid.Parse(idStr)
	if err != nil {
		response.BadRequest(c, "Invalid club ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User ID not found in context")
		return
	}
	userID := userIDVal.(uuid.UUID)

	requests, err := h.clubService.GetPendingJoinRequests(clubID, userID)
	if err != nil {
		if err == service.ErrNotAuthorized {
			response.Unauthorized(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, requests)
}

func (h *ClubHandler) ApproveJoinRequest(c *gin.Context) {
	clubIDStr := c.Param("id")
	_, err := uuid.Parse(clubIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid club ID")
		return
	}

	reqIDStr := c.Param("request_id")
	requestID, err := uuid.Parse(reqIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid request ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User ID not found in context")
		return
	}
	userID := userIDVal.(uuid.UUID)

	err = h.clubService.ApproveJoinRequest(requestID, userID)
	if err != nil {
		if err == service.ErrNotAuthorized {
			response.Unauthorized(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Join request approved successfully"})
}

func (h *ClubHandler) RejectJoinRequest(c *gin.Context) {
	clubIDStr := c.Param("id")
	_, err := uuid.Parse(clubIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid club ID")
		return
	}

	reqIDStr := c.Param("request_id")
	requestID, err := uuid.Parse(reqIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid request ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User ID not found in context")
		return
	}
	userID := userIDVal.(uuid.UUID)

	err = h.clubService.RejectJoinRequest(requestID, userID)
	if err != nil {
		if err == service.ErrNotAuthorized {
			response.Unauthorized(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Join request rejected successfully"})
}
type UpdateSettingsRequest struct {
	OwnerLabel      string `json:"owner_label"`
	AdminLabel      string `json:"admin_label"`
	LeadershipLabel string `json:"leadership_label"`
}

func (h *ClubHandler) UpdateSettings(c *gin.Context) {
	clubIDStr := c.Param("id")
	clubID, err := uuid.Parse(clubIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid club ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User ID not found in context")
		return
	}
	userID := userIDVal.(uuid.UUID)

	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	err = h.clubService.UpdateSettings(clubID, userID, req.OwnerLabel, req.AdminLabel, req.LeadershipLabel)
	if err != nil {
		if err == service.ErrNotAuthorized {
			response.Unauthorized(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Club settings updated successfully"})
}
