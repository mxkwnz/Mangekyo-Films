package handlers

import (
	"cinema-system/internal/models"
	"cinema-system/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReviewHandler struct {
	reviewService *services.ReviewService
}

func NewReviewHandler(reviewService *services.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

func (h *ReviewHandler) CreateReview(c *gin.Context) {
	userID := c.MustGet("userID").(primitive.ObjectID)

	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review.UserID = userID

	if err := h.reviewService.CreateReview(c.Request.Context(), &review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func (h *ReviewHandler) GetMovieReviews(c *gin.Context) {
	movieID, err := primitive.ObjectIDFromHex(c.Param("movieId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie ID"})
		return
	}

	reviews, err := h.reviewService.GetMovieReviews(c.Request.Context(), movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	reviewID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review ID"})
		return
	}

	userID := c.MustGet("userID").(primitive.ObjectID)
	role := c.GetHeader("X-User-Role")
	isAdmin := role == string(models.RoleAdmin)

	if err := h.reviewService.DeleteReview(c.Request.Context(), reviewID, userID, isAdmin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review deleted successfully"})
}
