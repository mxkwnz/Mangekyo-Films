package handlers

import (
	"cinema-system/internal/models"
	"cinema-system/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieHandler struct {
	movieService *services.MovieService
}

func NewMovieHandler(movieService *services.MovieService) *MovieHandler {
	return &MovieHandler{movieService: movieService}
}

func (h *MovieHandler) GetAllMovies(c *gin.Context) {
	movies, err := h.movieService.GetAllMovies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movies)
}

func (h *MovieHandler) GetMovie(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie ID"})
		return
	}

	movie, err := h.movieService.GetMovieByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

func (h *MovieHandler) CreateMovie(c *gin.Context) {
	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.movieService.CreateMovie(c.Request.Context(), &movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, movie)
}

func (h *MovieHandler) UpdateMovie(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie ID"})
		return
	}

	var movie models.Movie
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.movieService.UpdateMovie(c.Request.Context(), id, &movie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "movie updated successfully"})
}

func (h *MovieHandler) DeleteMovie(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie ID"})
		return
	}

	if err := h.movieService.DeleteMovie(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "movie deleted successfully"})
}
