package handler

import (
	"github.com/coolkit-org/coolkit/internal/model"
	"github.com/coolkit-org/coolkit/internal/repository"
	"github.com/coolkit-org/coolkit/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AnnouncementHandler struct {
	announcementRepo *repository.AnnouncementRepository
	clubRepo         *repository.ClubRepository
}

func NewAnnouncementHandler(announcementRepo *repository.AnnouncementRepository, clubRepo *repository.ClubRepository) *AnnouncementHandler {
	return &AnnouncementHandler{
		announcementRepo: announcementRepo,
		clubRepo:         clubRepo,
	}
}

type CreateAnnouncementRequest struct {
	Title    string `json:"title" binding:"required"`
	Content  string `json:"content" binding:"required"`
	Priority string `json:"priority"`
}

func (h *AnnouncementHandler) Create(c *gin.Context) {
	clubIDStr := c.Param("id")
	clubID, err := uuid.Parse(clubIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid club ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		response.Unauthorized(c, "Invalid user ID type in token")
		return
	}

	club, err := h.clubRepo.FindByID(clubID)
	if err != nil {
		response.NotFound(c, "Club not found")
		return
	}

	isAdmin := club.OwnerID == userID
	if !isAdmin {
		members, err := h.clubRepo.FindMembers(clubID)
		if err == nil {
			for _, m := range members {
				if m.UserID == userID && (m.Role == model.RoleAdmin || m.Role == model.RoleOwner) {
					isAdmin = true
					break
				}
			}
		}
	}

	if !isAdmin {
		response.Unauthorized(c, "Only club owners or admins can create announcements")
		return
	}

	var req CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	priority := req.Priority
	if priority == "" {
		priority = "normal"
	}

	announcement := &model.Announcement{
		ClubID:   clubID,
		AuthorID: userID,
		Title:    req.Title,
		Content:  req.Content,
		Priority: priority,
	}

	if err := h.announcementRepo.Create(announcement); err != nil {
		response.InternalError(c, "Failed to create announcement")
		return
	}

	response.Created(c, announcement)
}

func (h *AnnouncementHandler) List(c *gin.Context) {
	clubIDStr := c.Param("id")
	clubID, err := uuid.Parse(clubIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid club ID")
		return
	}

	announcements, err := h.announcementRepo.FindByClubID(clubID)
	if err != nil {
		response.InternalError(c, "Failed to list announcements")
		return
	}

	response.OK(c, announcements)
}

func (h *AnnouncementHandler) Delete(c *gin.Context) {
	announcementIDStr := c.Param("id")
	announcementID, err := uuid.Parse(announcementIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid announcement ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	userID, ok := userIDVal.(uuid.UUID)
	if !ok {
		response.Unauthorized(c, "Invalid user ID type in token")
		return
	}

	announcement, err := h.announcementRepo.FindByID(announcementID)
	if err != nil {
		response.NotFound(c, "Announcement not found")
		return
	}

	isAuthorized := announcement.AuthorID == userID
	if !isAuthorized {
		club, err := h.clubRepo.FindByID(announcement.ClubID)
		if err == nil {
			if club.OwnerID == userID {
				isAuthorized = true
			} else {
				members, err := h.clubRepo.FindMembers(announcement.ClubID)
				if err == nil {
					for _, m := range members {
						if m.UserID == userID && (m.Role == model.RoleAdmin || m.Role == model.RoleOwner) {
							isAuthorized = true
							break
						}
					}
				}
			}
		}
	}

	if !isAuthorized {
		response.Unauthorized(c, "Only the author or club admins can delete announcements")
		return
	}

	if err := h.announcementRepo.Delete(announcementID); err != nil {
		response.InternalError(c, "Failed to delete announcement")
		return
	}

	response.OK(c, map[string]string{"message": "Announcement deleted successfully"})
}
