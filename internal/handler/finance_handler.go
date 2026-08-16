package handler

import (
	"time"

	"github.com/coolkit-org/coolkit/internal/model"
	"github.com/coolkit-org/coolkit/internal/repository"
	"github.com/coolkit-org/coolkit/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FinanceHandler struct {
	financeRepo *repository.FinanceRepository
	eventRepo   *repository.EventRepository
	clubRepo    *repository.ClubRepository
}

func NewFinanceHandler(financeRepo *repository.FinanceRepository, eventRepo *repository.EventRepository, clubRepo *repository.ClubRepository) *FinanceHandler {
	return &FinanceHandler{
		financeRepo: financeRepo,
		eventRepo:   eventRepo,
		clubRepo:    clubRepo,
	}
}

type CreateFinanceRequest struct {
	Type        string    `json:"type" binding:"required"`
	Category    string    `json:"category" binding:"required"`
	Amount      float64   `json:"amount" binding:"required,min=0.01"`
	Description string    `json:"description"`
	Date        time.Time `json:"date" binding:"required"`
	ProofImage  string    `json:"proof_image"`
}

func (h *FinanceHandler) Create(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid event ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	userID := userIDVal.(uuid.UUID)

	event, err := h.eventRepo.FindByID(eventID)
	if err != nil {
		response.NotFound(c, "Event not found")
		return
	}

	members, err := h.clubRepo.FindMembers(event.ClubID)
	if err != nil {
		response.InternalError(c, "Failed to retrieve club members")
		return
	}

	var memberRole string
	isMember := false
	for _, m := range members {
		if m.UserID == userID {
			memberRole = m.Role
			isMember = true
			break
		}
	}

	if !isMember {
		response.Unauthorized(c, "You must be a club member to create finance entries")
		return
	}

	var req CreateFinanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Type != model.FinanceTypeIncome && req.Type != model.FinanceTypeExpense {
		response.BadRequest(c, "Invalid finance type")
		return
	}

	status := model.FinanceStatusPending
	if memberRole == model.RoleOwner || memberRole == model.RoleAdmin {
		status = model.FinanceStatusApproved
	}

	now := time.Now()
	entry := &model.FinanceEntry{
		EventID:     eventID,
		CreatedByID: userID,
		Type:        req.Type,
		Category:    req.Category,
		Amount:      req.Amount,
		Description: req.Description,
		Date:        req.Date,
		ProofImage:  req.ProofImage,
		Status:      status,
		ApprovedByID: func() *uuid.UUID {
			if status == model.FinanceStatusApproved {
				return &userID
			}
			return nil
		}(),
		ApprovedAt: func() *time.Time {
			if status == model.FinanceStatusApproved {
				return &now
			}
			return nil
		}(),
	}

	if err := h.financeRepo.Create(entry); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, entry)
}

func (h *FinanceHandler) List(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid event ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	userID := userIDVal.(uuid.UUID)

	event, err := h.eventRepo.FindByID(eventID)
	if err != nil {
		response.NotFound(c, "Event not found")
		return
	}

	members, err := h.clubRepo.FindMembers(event.ClubID)
	if err != nil {
		response.InternalError(c, "Failed to retrieve club members")
		return
	}

	var memberRole string
	isAdmin := false
	for _, m := range members {
		if m.UserID == userID {
			memberRole = m.Role
			if m.Role == model.RoleOwner || m.Role == model.RoleAdmin {
				isAdmin = true
			}
			break
		}
	}

	if memberRole == "" {
		response.Unauthorized(c, "You must be a club member to view finance entries")
		return
	}

	var entries []model.FinanceEntry
	if isAdmin {
		entries, err = h.financeRepo.FindByEventID(eventID)
	} else {
		entries, err = h.financeRepo.FindApprovedByEventID(eventID)
	}
	if err != nil {
		response.InternalError(c, "Failed to retrieve finance entries")
		return
	}

	totalIncome, totalExpense, err := h.financeRepo.GetSummary(eventID)
	if err != nil {
		response.InternalError(c, "Failed to calculate summary")
		return
	}

	summary := model.FinanceSummary{
		TotalIncome:  totalIncome,
		TotalExpense: totalExpense,
		Balance:      totalIncome - totalExpense,
		Entries:      entries,
	}

	if summary.Entries == nil {
		summary.Entries = []model.FinanceEntry{}
	}

	response.OK(c, summary)
}

func (h *FinanceHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid finance ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	userID := userIDVal.(uuid.UUID)

	entry, err := h.financeRepo.FindByID(id)
	if err != nil {
		response.NotFound(c, "Finance entry not found")
		return
	}

	if entry.CreatedByID != userID {
		event, err := h.eventRepo.FindByID(entry.EventID)
		if err != nil {
			response.NotFound(c, "Event not found")
			return
		}

		members, err := h.clubRepo.FindMembers(event.ClubID)
		if err != nil {
			response.InternalError(c, "Failed to retrieve club members")
			return
		}

		var memberRole string
		for _, m := range members {
			if m.UserID == userID {
				memberRole = m.Role
				break
			}
		}

		if memberRole != model.RoleOwner && memberRole != model.RoleAdmin {
			response.Unauthorized(c, "Only creator, club owner, or admin can delete finance entries")
			return
		}
	}

	if err := h.financeRepo.Delete(id); err != nil {
		response.InternalError(c, "Failed to delete finance entry")
		return
	}

	response.OK(c, map[string]string{"message": "Finance entry deleted successfully"})
}

type ApproveFinanceRequest struct {
	Approved bool   `json:"approved"`
	Reason   string `json:"reason"`
}

func (h *FinanceHandler) Approve(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid finance ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	userID := userIDVal.(uuid.UUID)

	entry, err := h.financeRepo.FindByID(id)
	if err != nil {
		response.NotFound(c, "Finance entry not found")
		return
	}

	event, err := h.eventRepo.FindByID(entry.EventID)
	if err != nil {
		response.NotFound(c, "Event not found")
		return
	}

	members, err := h.clubRepo.FindMembers(event.ClubID)
	if err != nil {
		response.InternalError(c, "Failed to retrieve club members")
		return
	}

	isAdmin := false
	for _, m := range members {
		if m.UserID == userID {
			if m.Role == model.RoleOwner || m.Role == model.RoleAdmin {
				isAdmin = true
			}
			break
		}
	}

	if !isAdmin {
		response.Unauthorized(c, "Only club owners and admins can approve finance entries")
		return
	}

	var req ApproveFinanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	now := time.Now()
	if req.Approved {
		err = h.financeRepo.UpdateStatus(id, model.FinanceStatusApproved, &userID, &now, "")
	} else {
		err = h.financeRepo.UpdateStatus(id, model.FinanceStatusRejected, &userID, &now, req.Reason)
	}

	if err != nil {
		response.InternalError(c, "Failed to update finance entry status")
		return
	}

	status := "approved"
	if !req.Approved {
		status = "rejected"
	}
	response.OK(c, map[string]string{"message": "Finance entry " + status + " successfully"})
}
