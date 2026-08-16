package handler

import (
"github.com/gin-gonic/gin"
"github.com/google/uuid"
"net/http"
"github.com/coolkit-org/coolkit/internal/model"
"github.com/coolkit-org/coolkit/internal/service"
"github.com/coolkit-org/coolkit/pkg/response"
)

type RSVPHandler struct {
rsvpService service.RSVPService
}

func NewRSVPHandler(rsvpService service.RSVPService) *RSVPHandler {
return &RSVPHandler{rsvpService: rsvpService}
}

type CreateRSVPRequest struct {
Status string `json:"status" binding:"required"`
}

func (h *RSVPHandler) CreateOrUpdateRSVP(c *gin.Context) {
eventIDStr := c.Param("event_id")
eventID, err := uuid.Parse(eventIDStr)
if err != nil {
response.BadRequest(c, "Invalid event ID")
return
}

var req CreateRSVPRequest
if err := c.ShouldBindJSON(&req); err != nil {
response.BadRequest(c, err.Error())
return
}

userID, exists := c.Get("user_id")
if !exists {
response.Unauthorized(c, "User not authenticated")
return
}

rsvp := &model.RSVP{
EventID: eventID,
UserID:  userID.(uuid.UUID),
Status:  model.RSVPStatus(req.Status),
}

if err := h.rsvpService.CreateOrUpdateRSVP(c.Request.Context(), rsvp); err != nil {
response.Error(c, http.StatusInternalServerError, err.Error())
return
}

response.Success(c, http.StatusCreated, rsvp)
}

func (h *RSVPHandler) GetUserRSVP(c *gin.Context) {
eventIDStr := c.Param("event_id")
eventID, err := uuid.Parse(eventIDStr)
if err != nil {
response.BadRequest(c, "Invalid event ID")
return
}

userID, exists := c.Get("user_id")
if !exists {
response.Unauthorized(c, "User not authenticated")
return
}

rsvp, err := h.rsvpService.GetUserRSVP(c.Request.Context(), eventID, userID.(uuid.UUID))
if err != nil {
response.NotFound(c, "No RSVP found for this event")
return
}

response.Success(c, http.StatusOK, rsvp)
}

func (h *RSVPHandler) GetEventRSVPs(c *gin.Context) {
eventIDStr := c.Param("event_id")
eventID, err := uuid.Parse(eventIDStr)
if err != nil {
response.BadRequest(c, "Invalid event ID")
return
}

rsvps, err := h.rsvpService.GetEventRSVPs(c.Request.Context(), eventID)
if err != nil {
response.Error(c, http.StatusInternalServerError, "Failed to fetch RSVPs")
return
}

response.Success(c, http.StatusOK, rsvps)
}

func (h *RSVPHandler) GetRSVPCounts(c *gin.Context) {
eventIDStr := c.Param("event_id")
eventID, err := uuid.Parse(eventIDStr)
if err != nil {
response.BadRequest(c, "Invalid event ID")
return
}

counts, err := h.rsvpService.GetRSVPCounts(c.Request.Context(), eventID)
if err != nil {
response.Error(c, http.StatusInternalServerError, "Failed to fetch RSVP counts")
return
}

response.Success(c, http.StatusOK, counts)
}

func (h *RSVPHandler) CheckInUser(c *gin.Context) {
rsvpIDStr := c.Param("rsvp_id")
rsvpID, err := uuid.Parse(rsvpIDStr)
if err != nil {
response.BadRequest(c, "Invalid RSVP ID")
return
}

if err := h.rsvpService.CheckInUser(c.Request.Context(), rsvpID); err != nil {
response.Error(c, http.StatusInternalServerError, err.Error())
return
}

response.Success(c, http.StatusOK, nil)
}
