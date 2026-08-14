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

	club, err := h.clubService.CreateClub(req.Name, req.Description, userID)
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

func (h *ClubHandler) Join(c *gin.Context) {
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

	err = h.clubService.JoinClub(clubID, userID)
	if err != nil {
		if errors.Is(err, service.ErrAlreadyMember) {
			response.Error(c, 409, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Successfully joined club"})
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
