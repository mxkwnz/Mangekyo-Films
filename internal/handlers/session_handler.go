package handlers

import (
	"cinema-system/internal/models"
	"cinema-system/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionHandler struct {
	sessionService *services.SessionService
}

func NewSessionHandler(sessionService *services.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

type CreateSessionRequest struct {
	MovieID   string  `json:"movie_id" binding:"required"`
	HallID    string  `json:"hall_id" binding:"required"`
	StartTime string  `json:"start_time" binding:"required"`
	Price     float64 `json:"price"`
}

func (h *SessionHandler) CreateSession(c *gin.Context) {
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	movieID, err := primitive.ObjectIDFromHex(req.MovieID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID format"})
		return
	}

	hallID, err := primitive.ObjectIDFromHex(req.HallID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hall ID format"})
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format (expected ISO-8601)"})
		return
	}

	session := models.Session{
		MovieID:   movieID,
		HallID:    hallID,
		StartTime: startTime,
		Price:     req.Price,
	}

	if err := h.sessionService.CreateSession(c.Request.Context(), &session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, session)
}

func (h *SessionHandler) GetUpcomingSessions(c *gin.Context) {
	sessions, err := h.sessionService.GetUpcomingSessions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func (h *SessionHandler) GetUpcomingMovieIDs(c *gin.Context) {
	ids, err := h.sessionService.GetUpcomingMovieIDs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ids)
}

func (h *SessionHandler) GetMovieSessions(c *gin.Context) {
	movieID, err := primitive.ObjectIDFromHex(c.Param("movieId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie ID"})
		return
	}

	sessions, err := h.sessionService.GetSessionsByMovie(c.Request.Context(), movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sessions)
}

func (h *SessionHandler) GetSession(c *gin.Context) {
	sessionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}
	session, err := h.sessionService.GetSessionByID(c.Request.Context(), sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}
	c.JSON(http.StatusOK, session)
}

func (h *SessionHandler) UpdateSession(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request: " + err.Error()})
		return
	}

	movieID, err := primitive.ObjectIDFromHex(req.MovieID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID format"})
		return
	}

	hallID, err := primitive.ObjectIDFromHex(req.HallID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid hall ID format"})
		return
	}

	startTime, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start time format (expected ISO-8601)"})
		return
	}

	session := models.Session{
		MovieID:   movieID,
		HallID:    hallID,
		StartTime: startTime,
		Price:     req.Price,
	}

	if err := h.sessionService.UpdateSession(c.Request.Context(), id, &session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session updated successfully"})
}

func (h *SessionHandler) DeleteSession(c *gin.Context) {
	sessionID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	if err := h.sessionService.DeleteSession(c.Request.Context(), sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session deleted successfully"})
}
