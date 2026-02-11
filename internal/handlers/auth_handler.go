package handlers

import (
	"cinema-system/internal/models"
	"cinema-system/internal/services"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	go_playground_validator "github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.UserRegistration
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	user, token, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.AuthResponse{
		User:    user,
		Token:   token,
		Message: "registration successful",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.UserLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": formatValidationError(err)})
		return
	}

	user, token, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.AuthResponse{
		User:    user,
		Token:   token,
		Message: "login successful",
	})
}

func formatValidationError(err error) string {
	var ve go_playground_validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			switch fe.Field() {
			case "Email":
				return "Please enter a valid email address"
			case "Password":
				if fe.Tag() == "min" {
					return "Password must be at least 6 characters long"
				}
				return "Password is required"
			case "FirstName":
				return "First name must be at least 2 characters long"
			case "LastName":
				return "Last name must be at least 2 characters long"
			case "PhoneNumber":
				return "Phone number is required"
			}
		}
	}
	return "Invalid input data"
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.authService.GetUserByID(c.Request.Context(), userID.(primitive.ObjectID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID := c.MustGet("userID").(primitive.ObjectID)

	var req struct {
		FirstName   string `json:"first_name" binding:"required,min=2"`
		LastName    string `json:"last_name" binding:"required,min=2"`
		Email       string `json:"email" binding:"required,email"`
		PhoneNumber string `json:"phone_number"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.UpdateProfile(c.Request.Context(), userID, req.FirstName, req.LastName, req.Email, req.PhoneNumber); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}
