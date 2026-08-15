package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/coolkit-org/coolkit/internal/model"
	"github.com/coolkit-org/coolkit/internal/repository"
	"github.com/coolkit-org/coolkit/pkg/response"
)

type HierarchyHandler struct {
	clubRepo      *repository.ClubRepository
	hierarchyRepo *repository.HierarchyRepository
	domainRepo    *repository.DomainRepository
}

func NewHierarchyHandler(clubRepo *repository.ClubRepository, hierarchyRepo *repository.HierarchyRepository, domainRepo *repository.DomainRepository) *HierarchyHandler {
	return &HierarchyHandler{
		clubRepo:      clubRepo,
		hierarchyRepo: hierarchyRepo,
		domainRepo:    domainRepo,
	}
}

func isClubAdmin(clubRepo *repository.ClubRepository, clubID, userID uuid.UUID) (bool, error) {
	members, err := clubRepo.FindMembers(clubID)
	if err != nil {
		return false, err
	}
	for _, m := range members {
		if m.UserID == userID && (m.Role == model.RoleOwner || m.Role == model.RoleAdmin) {
			return true, nil
		}
	}
	return false, nil
}

type SetHierarchyRequest struct {
	Levels []struct {
		Name     string `json:"name" binding:"required"`
		Position int    `json:"position" binding:"required"`
	} `json:"levels" binding:"required"`
}

func (h *HierarchyHandler) SetHierarchy(c *gin.Context) {
	clubID, err := uuid.Parse(c.Param("id"))
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

	isAdmin, err := isClubAdmin(h.clubRepo, clubID, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if !isAdmin {
		response.Forbidden(c, "Only owner or admin can set hierarchy")
		return
	}

	var req SetHierarchyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var levels []model.HierarchyLevel
	for _, l := range req.Levels {
		levels = append(levels, model.HierarchyLevel{
			ClubID:   clubID,
			Name:     l.Name,
			Position: l.Position,
		})
	}

	if err := h.hierarchyRepo.BulkReplace(clubID, levels); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	newLevels, err := h.hierarchyRepo.FindByClubID(clubID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, newLevels)
}

func (h *HierarchyHandler) GetHierarchy(c *gin.Context) {
	clubID, err := uuid.Parse(c.Param("id"))
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

	isMember, err := h.clubRepo.IsMember(clubID, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if !isMember {
		response.Forbidden(c, "Only members can view hierarchy")
		return
	}

	levels, err := h.hierarchyRepo.FindByClubID(clubID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, levels)
}

type CreateDomainRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (h *HierarchyHandler) CreateDomain(c *gin.Context) {
	clubID, err := uuid.Parse(c.Param("id"))
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

	isAdmin, err := isClubAdmin(h.clubRepo, clubID, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if !isAdmin {
		response.Forbidden(c, "Only owner or admin can create domains")
		return
	}

	var req CreateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	domain := &model.Domain{
		ClubID:      clubID,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.domainRepo.Create(domain); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, domain)
}

func (h *HierarchyHandler) ListDomains(c *gin.Context) {
	clubID, err := uuid.Parse(c.Param("id"))
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

	isMember, err := h.clubRepo.IsMember(clubID, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if !isMember {
		response.Forbidden(c, "Only members can view domains")
		return
	}

	domains, err := h.domainRepo.FindByClubID(clubID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, domains)
}

func (h *HierarchyHandler) UpdateDomain(c *gin.Context) {
	clubID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid club ID")
		return
	}

	domainID, err := uuid.Parse(c.Param("domain_id"))
	if err != nil {
		response.BadRequest(c, "Invalid domain ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User ID not found in context")
		return
	}
	userID := userIDVal.(uuid.UUID)

	isAdmin, err := isClubAdmin(h.clubRepo, clubID, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if !isAdmin {
		response.Forbidden(c, "Only owner or admin can update domains")
		return
	}

	var req CreateDomainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	domain, err := h.domainRepo.FindByID(domainID)
	if err != nil {
		response.NotFound(c, "Domain not found")
		return
	}

	if domain.ClubID != clubID {
		response.BadRequest(c, "Domain does not belong to this club")
		return
	}

	domain.Name = req.Name
	domain.Description = req.Description

	if err := h.domainRepo.Update(domain); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, domain)
}

func (h *HierarchyHandler) DeleteDomain(c *gin.Context) {
	clubID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid club ID")
		return
	}

	domainID, err := uuid.Parse(c.Param("domain_id"))
	if err != nil {
		response.BadRequest(c, "Invalid domain ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User ID not found in context")
		return
	}
	userID := userIDVal.(uuid.UUID)

	isAdmin, err := isClubAdmin(h.clubRepo, clubID, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if !isAdmin {
		response.Forbidden(c, "Only owner or admin can delete domains")
		return
	}

	domain, err := h.domainRepo.FindByID(domainID)
	if err != nil {
		response.NotFound(c, "Domain not found")
		return
	}

	if domain.ClubID != clubID {
		response.BadRequest(c, "Domain does not belong to this club")
		return
	}

	// Update members in this domain to null domain
	members, err := h.clubRepo.FindMembers(clubID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	for _, m := range members {
		if m.DomainID != nil && *m.DomainID == domainID {
			if err := h.clubRepo.UpdateMemberOrganization(clubID, m.UserID, nil, m.HierarchyLevelID); err != nil {
				response.InternalError(c, err.Error())
				return
			}
		}
	}

	if err := h.domainRepo.Delete(domainID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Domain deleted successfully"})
}

type AssignOrganizationRequest struct {
	DomainID         *string `json:"domain_id"`
	HierarchyLevelID *string `json:"hierarchy_level_id"`
}

func (h *HierarchyHandler) AssignMemberOrganization(c *gin.Context) {
	clubID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid club ID")
		return
	}

	targetUserID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		response.BadRequest(c, "Invalid target user ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User ID not found in context")
		return
	}
	userID := userIDVal.(uuid.UUID)

	isAdmin, err := isClubAdmin(h.clubRepo, clubID, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if !isAdmin {
		response.Forbidden(c, "Only owner or admin can assign organization")
		return
	}

	var req AssignOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var domainID *uuid.UUID
	if req.DomainID != nil && *req.DomainID != "" {
		id, err := uuid.Parse(*req.DomainID)
		if err != nil {
			response.BadRequest(c, "Invalid domain ID format")
			return
		}
		domainID = &id
	}

	var hierarchyLevelID *uuid.UUID
	if req.HierarchyLevelID != nil && *req.HierarchyLevelID != "" {
		id, err := uuid.Parse(*req.HierarchyLevelID)
		if err != nil {
			response.BadRequest(c, "Invalid hierarchy level ID format")
			return
		}
		hierarchyLevelID = &id
	}

	if err := h.clubRepo.UpdateMemberOrganization(clubID, targetUserID, domainID, hierarchyLevelID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Organization assigned successfully"})
}

func (h *HierarchyHandler) GetOrganizationTree(c *gin.Context) {
	clubID, err := uuid.Parse(c.Param("id"))
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

	isMember, err := h.clubRepo.IsMember(clubID, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	if !isMember {
		response.Forbidden(c, "Only members can view organization tree")
		return
	}

	members, err := h.clubRepo.FindMembers(clubID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	levels, err := h.hierarchyRepo.FindByClubID(clubID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	domains, err := h.domainRepo.FindByClubID(clubID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{
		"members":          members,
		"hierarchy_levels": levels,
		"domains":          domains,
	})
}
