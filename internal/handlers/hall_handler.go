package handlers

import (
	"cinema-system/internal/models"
	"cinema-system/internal/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HallHandler struct {
	hallRepo *repositories.HallRepository
}

func NewHallHandler(hallRepo *repositories.HallRepository) *HallHandler {
	return &HallHandler{hallRepo: hallRepo}
}

func (h *HallHandler) CreateHall(c *gin.Context) {
	var hall models.Hall
	if err := c.ShouldBindJSON(&hall); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.hallRepo.Create(c.Request.Context(), &hall); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, hall)
}

func (h *HallHandler) GetAllHalls(c *gin.Context) {
	halls, err := h.hallRepo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, halls)
}

func (h *HallHandler) GetHall(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hall ID"})
		return
	}
	hall, err := h.hallRepo.FindByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "hall not found"})
		return
	}
	c.JSON(http.StatusOK, hall)
}

func (h *HallHandler) UpdateHall(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hall ID"})
		return
	}

	var hall models.Hall
	if err := c.ShouldBindJSON(&hall); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.hallRepo.Update(c.Request.Context(), id, &hall); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "hall updated successfully"})
}

func (h *HallHandler) DeleteHall(c *gin.Context) {
	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid hall ID"})
		return
	}

	if err := h.hallRepo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "hall deleted successfully"})
}
