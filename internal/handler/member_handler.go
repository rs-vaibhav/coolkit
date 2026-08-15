package handler

import (
	"github.com/coolkit-org/coolkit/internal/service"
	"github.com/coolkit-org/coolkit/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MemberHandler struct {
	clubService *service.ClubService
}

func NewMemberHandler(clubService *service.ClubService) *MemberHandler {
	return &MemberHandler{clubService: clubService}
}

type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required"`
}

func (h *MemberHandler) UpdateRole(c *gin.Context) {
	clubID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid club ID")
		return
	}

	targetUserID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	requesterID := userIDVal.(uuid.UUID)

	var req UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	err = h.clubService.UpdateMemberRole(clubID, requesterID, targetUserID, req.Role)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Role updated successfully"})
}

func (h *MemberHandler) Remove(c *gin.Context) {
	clubID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid club ID")
		return
	}

	targetUserID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		response.BadRequest(c, "Invalid user ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	requesterID := userIDVal.(uuid.UUID)

	err = h.clubService.RemoveMember(clubID, requesterID, targetUserID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Member removed successfully"})
}

func (h *MemberHandler) Leave(c *gin.Context) {
	clubID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid club ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	userID := userIDVal.(uuid.UUID)

	err = h.clubService.LeaveClub(clubID, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "You have left the club"})
}
