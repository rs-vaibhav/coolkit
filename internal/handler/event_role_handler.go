package handler

import (
	"github.com/coolkit-org/coolkit/internal/service"
	"github.com/coolkit-org/coolkit/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventRoleHandler struct {
	eventService *service.EventService
}

func NewEventRoleHandler(eventService *service.EventService) *EventRoleHandler {
	return &EventRoleHandler{eventService: eventService}
}

func (h *EventRoleHandler) GetEventDetails(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid event ID")
		return
	}

	event, err := h.eventService.GetEvent(eventID)
	if err != nil {
		response.NotFound(c, "Event not found")
		return
	}

	response.OK(c, event)
}

func (h *EventRoleHandler) GetRoles(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid event ID")
		return
	}

	roles, err := h.eventService.GetEventRoles(eventID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, roles)
}

type AssignRoleRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required"`
	RoleName string    `json:"role_name" binding:"required"`
}

func (h *EventRoleHandler) AssignRole(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid event ID")
		return
	}

	assignerIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	assignerID := assignerIDVal.(uuid.UUID)

	var req AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	role, err := h.eventService.AssignEventRole(eventID, assignerID, req.UserID, req.RoleName)
	if err != nil {
		if err == service.ErrUnauthorized {
			response.Unauthorized(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, role)
}

func (h *EventRoleHandler) RemoveRole(c *gin.Context) {
	roleID, err := uuid.Parse(c.Param("role_id"))
	if err != nil {
		response.BadRequest(c, "Invalid role ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	userID := userIDVal.(uuid.UUID)

	err = h.eventService.RemoveEventRole(roleID, userID)
	if err != nil {
		if err == service.ErrUnauthorized {
			response.Unauthorized(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Role removed successfully"})
}
