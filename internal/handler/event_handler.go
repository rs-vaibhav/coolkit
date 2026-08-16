package handler

import (
"time"

"github.com/coolkit-org/coolkit/internal/service"
"github.com/coolkit-org/coolkit/pkg/response"
"github.com/gin-gonic/gin"
"github.com/google/uuid"
)

type EventHandler struct {
eventService *service.EventService
}

func NewEventHandler(eventService *service.EventService) *EventHandler {
return &EventHandler{eventService: eventService}
}

type CreateEventRequest struct {
Title       string    `json:"title" binding:"required"`
Description string    `json:"description"`
Date        time.Time `json:"date" binding:"required"`
Location    string    `json:"location"`
}

func (h *EventHandler) Create(c *gin.Context) {
clubIDStr := c.Param("id")
clubID, err := uuid.Parse(clubIDStr)
if err != nil {
response.BadRequest(c, "Invalid club ID")
return
}

userIDStr, exists := c.Get("user_id")
if !exists {
response.Unauthorized(c, "Unauthorized")
return
}

userID, ok := userIDStr.(uuid.UUID)
if !ok {
response.Unauthorized(c, "Invalid user ID type in token")
return
}

var req CreateEventRequest
if err := c.ShouldBindJSON(&req); err != nil {
response.BadRequest(c, err.Error())
return
}

event, err := h.eventService.CreateEvent(clubID, userID, req.Title, req.Description, req.Date, req.Location)
if err != nil {
if err == service.ErrUnauthorized {
response.Unauthorized(c, err.Error())
return
}
response.InternalError(c, err.Error())
return
}

response.Created(c, event)
}

func (h *EventHandler) List(c *gin.Context) {
clubIDStr := c.Param("id")
clubID, err := uuid.Parse(clubIDStr)
if err != nil {
response.BadRequest(c, "Invalid club ID")
return
}

events, err := h.eventService.GetClubEvents(clubID)
if err != nil {
response.InternalError(c, err.Error())
return
}

response.OK(c, events)
}

type UpdateEventRequest struct {
Title             string    `json:"title"`
Description       string    `json:"description"`
Date              time.Time `json:"date"`
Location          string    `json:"location"`
FormbricksEnvID   string    `json:"formbricks_env_id,omitempty"`
FormbricksSurveyID string   `json:"formbricks_survey_id,omitempty"`
}

func (h *EventHandler) Update(c *gin.Context) {
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

var req UpdateEventRequest
if err := c.ShouldBindJSON(&req); err != nil {
response.BadRequest(c, err.Error())
return
}

event, err := h.eventService.UpdateEvent(eventID, userID, req.Title, req.Description, req.Date, req.Location)
if err != nil {
if err == service.ErrUnauthorized {
response.Unauthorized(c, err.Error())
return
}
response.InternalError(c, err.Error())
return
}

response.OK(c, event)
}

func (h *EventHandler) Delete(c *gin.Context) {
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

err = h.eventService.DeleteEvent(eventID, userID)
if err != nil {
if err == service.ErrUnauthorized {
response.Unauthorized(c, err.Error())
return
}
response.InternalError(c, err.Error())
return
}

response.OK(c, gin.H{"message": "Event deleted successfully"})
}

// UpdateEventWithFormbricks updates event including Formbricks configuration
func (h *EventHandler) UpdateEventWithFormbricks(c *gin.Context) {
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

var req UpdateEventRequest
if err := c.ShouldBindJSON(&req); err != nil {
response.BadRequest(c, err.Error())
return
}

event, err := h.eventService.UpdateEventFormbricks(eventID, userID, req.FormbricksEnvID, req.FormbricksSurveyID)
if err != nil {
if err == service.ErrUnauthorized {
response.Unauthorized(c, err.Error())
return
}
response.InternalError(c, err.Error())
return
}

response.OK(c, event)
}

// CreateFormbricksSurvey creates a new Formbricks survey for the event
func (h *EventHandler) CreateFormbricksSurvey(c *gin.Context) {
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

type Request struct {
SurveyName        string `json:"survey_name" binding:"required"`
SurveyDescription string `json:"survey_description"`
}

var req Request
if err := c.ShouldBindJSON(&req); err != nil {
response.BadRequest(c, err.Error())
return
}

event, err := h.eventService.CreateFormbricksSurvey(eventID, userID, req.SurveyName, req.SurveyDescription)
if err != nil {
if err == service.ErrUnauthorized {
response.Unauthorized(c, err.Error())
return
}
response.InternalError(c, err.Error())
return
}

response.OK(c, gin.H{
"message":           "Survey created successfully",
"survey_id":         event.FormbricksSurveyID,
"formbricks_env_id": event.FormbricksEnvID,
})
}
