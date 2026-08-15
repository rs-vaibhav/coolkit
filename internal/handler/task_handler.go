package handler

import (
	"github.com/coolkit-org/coolkit/internal/model"
	"github.com/coolkit-org/coolkit/internal/repository"
	"github.com/coolkit-org/coolkit/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"time"
)

type TaskHandler struct {
	taskRepo  *repository.TaskRepository
	eventRepo *repository.EventRepository
	clubRepo  *repository.ClubRepository
}

func NewTaskHandler(taskRepo *repository.TaskRepository, eventRepo *repository.EventRepository, clubRepo *repository.ClubRepository) *TaskHandler {
	return &TaskHandler{
		taskRepo:  taskRepo,
		eventRepo: eventRepo,
		clubRepo:  clubRepo,
	}
}

type CreateTaskRequest struct {
	Title        string     `json:"title" binding:"required"`
	Description  string     `json:"description"`
	AssignedToID uuid.UUID  `json:"assigned_to_id" binding:"required"`
	DueDate      *time.Time `json:"due_date"`
}

func (h *TaskHandler) Create(c *gin.Context) {
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
		response.InternalError(c, "Failed to fetch club members")
		return
	}

	isAuthorized := false
	for _, m := range members {
		if m.UserID == userID && (m.Role == model.RoleOwner || m.Role == model.RoleAdmin) {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		response.Unauthorized(c, "Only club owner or admin can create tasks")
		return
	}

	var req CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	task := &model.Task{
		EventID:      eventID,
		AssignedToID: req.AssignedToID,
		Title:        req.Title,
		Description:  req.Description,
		DueDate:      req.DueDate,
		Status:       model.TaskStatusTodo,
	}

	if err := h.taskRepo.Create(task); err != nil {
		response.InternalError(c, "Failed to create task")
		return
	}

	response.Created(c, task)
}

func (h *TaskHandler) List(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid event ID")
		return
	}

	tasks, err := h.taskRepo.FindByEventID(eventID)
	if err != nil {
		response.InternalError(c, "Failed to fetch tasks")
		return
	}

	response.OK(c, tasks)
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

func (h *TaskHandler) UpdateStatus(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid task ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	userID := userIDVal.(uuid.UUID)

	task, err := h.taskRepo.FindByID(taskID)
	if err != nil {
		response.NotFound(c, "Task not found")
		return
	}

	event, err := h.eventRepo.FindByID(task.EventID)
	if err != nil {
		response.InternalError(c, "Failed to fetch event details")
		return
	}

	members, err := h.clubRepo.FindMembers(event.ClubID)
	if err != nil {
		response.InternalError(c, "Failed to fetch club members")
		return
	}

	isAuthorized := false
	if task.AssignedToID == userID {
		isAuthorized = true
	} else {
		for _, m := range members {
			if m.UserID == userID && (m.Role == model.RoleOwner || m.Role == model.RoleAdmin) {
				isAuthorized = true
				break
			}
		}
	}

	if !isAuthorized {
		response.Unauthorized(c, "Only assignee or club admin can update task status")
		return
	}

	var req UpdateTaskStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.Status != model.TaskStatusTodo && req.Status != model.TaskStatusInProgress && req.Status != model.TaskStatusDone {
		response.BadRequest(c, "Invalid status. Must be todo, in_progress, or done")
		return
	}

	if err := h.taskRepo.UpdateStatus(taskID, req.Status); err != nil {
		response.InternalError(c, "Failed to update task status")
		return
	}

	task.Status = req.Status
	response.OK(c, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid task ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "Unauthorized")
		return
	}
	userID := userIDVal.(uuid.UUID)

	task, err := h.taskRepo.FindByID(taskID)
	if err != nil {
		response.NotFound(c, "Task not found")
		return
	}

	event, err := h.eventRepo.FindByID(task.EventID)
	if err != nil {
		response.InternalError(c, "Failed to fetch event details")
		return
	}

	members, err := h.clubRepo.FindMembers(event.ClubID)
	if err != nil {
		response.InternalError(c, "Failed to fetch club members")
		return
	}

	isAuthorized := false
	for _, m := range members {
		if m.UserID == userID && (m.Role == model.RoleOwner || m.Role == model.RoleAdmin) {
			isAuthorized = true
			break
		}
	}

	if !isAuthorized {
		response.Unauthorized(c, "Only club owner or admin can delete tasks")
		return
	}

	if err := h.taskRepo.Delete(taskID); err != nil {
		response.InternalError(c, "Failed to delete task")
		return
	}

	response.OK(c, map[string]string{"message": "Task deleted successfully"})
}
