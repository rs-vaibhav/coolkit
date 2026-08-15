package handler

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/coolkit-org/coolkit/internal/service"
	"github.com/coolkit-org/coolkit/pkg/response"
)

type BookingHandler struct {
	bookingService *service.BookingService
}

func NewBookingHandler(bookingService *service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

type AddResourceRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity" binding:"required,min=1"`
}

func (h *BookingHandler) AddResource(c *gin.Context) {
	clubIDStr := c.Param("id")
	clubID, err := uuid.Parse(clubIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid club ID")
		return
	}

	var req AddResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User ID not found in context")
		return
	}
	userID := userIDVal.(uuid.UUID)

	res, err := h.bookingService.AddResource(req.Name, req.Description, req.Quantity, clubID, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotAuthorized) {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, res)
}

func (h *BookingHandler) ListResources(c *gin.Context) {
	clubIDStr := c.Param("id")
	clubID, err := uuid.Parse(clubIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid club ID")
		return
	}

	resources, err := h.bookingService.ListResources(clubID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, resources)
}

func (h *BookingHandler) DeleteResource(c *gin.Context) {
	resIDStr := c.Param("resource_id")
	resourceID, err := uuid.Parse(resIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid resource ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User ID not found in context")
		return
	}
	userID := userIDVal.(uuid.UUID)

	err = h.bookingService.DeleteResource(resourceID, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotAuthorized) {
			response.Forbidden(c, err.Error())
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Resource deleted successfully"})
}

type BookResourceRequest struct {
	ClubID    string    `json:"club_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required"`
	Purpose   string    `json:"purpose"`
}

func (h *BookingHandler) BookResource(c *gin.Context) {
	resIDStr := c.Param("resource_id")
	resourceID, err := uuid.Parse(resIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid resource ID")
		return
	}

	var req BookResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	clubID, err := uuid.Parse(req.ClubID)
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

	booking, err := h.bookingService.BookResource(resourceID, clubID, userID, req.StartTime, req.EndTime, req.Purpose)
	if err != nil {
		if errors.Is(err, service.ErrNotAuthorized) {
			response.Forbidden(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.Created(c, booking)
}

func (h *BookingHandler) ListBookings(c *gin.Context) {
	clubIDStr := c.Param("id")
	clubID, err := uuid.Parse(clubIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid club ID")
		return
	}

	bookings, err := h.bookingService.ListClubBookings(clubID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.OK(c, bookings)
}

func (h *BookingHandler) ApproveBooking(c *gin.Context) {
	bookingIDStr := c.Param("booking_id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid booking ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User ID not found in context")
		return
	}
	userID := userIDVal.(uuid.UUID)

	err = h.bookingService.ApproveBooking(bookingID, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotAuthorized) {
			response.Forbidden(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Booking approved successfully"})
}

func (h *BookingHandler) RejectBooking(c *gin.Context) {
	bookingIDStr := c.Param("booking_id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid booking ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User ID not found in context")
		return
	}
	userID := userIDVal.(uuid.UUID)

	err = h.bookingService.RejectBooking(bookingID, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotAuthorized) {
			response.Forbidden(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Booking rejected successfully"})
}

func (h *BookingHandler) CancelBooking(c *gin.Context) {
	bookingIDStr := c.Param("booking_id")
	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		response.BadRequest(c, "Invalid booking ID")
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "User ID not found in context")
		return
	}
	userID := userIDVal.(uuid.UUID)

	err = h.bookingService.CancelBooking(bookingID, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotAuthorized) {
			response.Forbidden(c, err.Error())
			return
		}
		response.BadRequest(c, err.Error())
		return
	}

	response.OK(c, gin.H{"message": "Booking cancelled/deleted successfully"})
}
