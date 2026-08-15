package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/coolkit-org/coolkit/internal/service"
	"github.com/coolkit-org/coolkit/pkg/response"
)

type AnalyticsHandler struct {
	analyticsService *service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

func (h *AnalyticsHandler) GetClubAnalytics(c *gin.Context) {
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

	stats, err := h.analyticsService.GetClubAnalytics(clubID, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotAuthorized) {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, stats)
}
