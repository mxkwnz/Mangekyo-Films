package handlers

import (
	"cinema-system/internal/models"
	"cinema-system/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingHandler struct {
	bookingService *services.BookingService
}

func NewBookingHandler(bookingService *services.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService}
}

type BookTicketSeat struct {
	RowNumber  int               `json:"row_number" binding:"required,min=1"`
	SeatNumber int               `json:"seat_number" binding:"required,min=1"`
	Type       models.TicketType `json:"type" binding:"required"`
}

type BookTicketsRequest struct {
	SessionID string           `json:"session_id" binding:"required"`
	Seats     []BookTicketSeat `json:"seats" binding:"required,dive"`
}

func (h *BookingHandler) BookTickets(c *gin.Context) {
	userID := c.MustGet("userID").(primitive.ObjectID)

	var req BookTicketsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionID, err := primitive.ObjectIDFromHex(req.SessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	var seatRequests []services.SeatBookingRequest
	for _, s := range req.Seats {
		seatRequests = append(seatRequests, services.SeatBookingRequest{
			RowNumber:  s.RowNumber,
			SeatNumber: s.SeatNumber,
			Type:       s.Type,
		})
	}

	tickets, err := h.bookingService.BookTickets(c.Request.Context(), userID, sessionID, seatRequests)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tickets)
}

func (h *BookingHandler) CancelTicket(c *gin.Context) {
	userID := c.MustGet("userID").(primitive.ObjectID)
	ticketID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket ID"})
		return
	}

	if err := h.bookingService.CancelTicket(c.Request.Context(), ticketID, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ticket cancelled successfully"})
}

func (h *BookingHandler) GetMyTickets(c *gin.Context) {
	userID := c.MustGet("userID").(primitive.ObjectID)

	tickets, err := h.bookingService.GetUserTickets(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func (h *BookingHandler) GetSessionTickets(c *gin.Context) {
	sessionID, err := primitive.ObjectIDFromHex(c.Param("sessionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	tickets, err := h.bookingService.GetSessionTickets(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func (h *BookingHandler) GetAllBookings(c *gin.Context) {
	tickets, err := h.bookingService.GetAllBookings(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tickets)
}

func (h *BookingHandler) GetSessionBookedSeats(c *gin.Context) {
	sessionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}
	seats, err := h.bookingService.GetSessionBookedSeats(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, seats)
}
