package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/coolkit-org/coolkit/internal/database"
	"github.com/coolkit-org/coolkit/pkg/response"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

func (h *HealthHandler) Health(c *gin.Context) {
	response.OK(c, gin.H{
		"status":  "ok",
		"version": "0.1.0",
	})
}

func (h *HealthHandler) DBHealth(c *gin.Context) {
	err := database.HealthCheck(h.db)
	if err != nil {
		response.InternalError(c, "database health check failed")
		return
	}
	response.OK(c, gin.H{
		"status": "ok",
	})
}
