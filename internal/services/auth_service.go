package services

import (
	"context"
	"errors"
	"time"

	"cinema-system/internal/models"
	"cinema-system/internal/repositories"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user with this email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthService struct {
	userRepo   *repositories.UserRepository
	reviewRepo *repositories.ReviewRepository
}

func NewAuthService(userRepo *repositories.UserRepository, reviewRepo *repositories.ReviewRepository) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		reviewRepo: reviewRepo,
	}
}

func (s *AuthService) Register(ctx context.Context, req models.UserRegistration) (*models.User, string, error) {
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, "", ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		PasswordHash: string(hashedPassword),
		Role:         models.RoleUser,
		Balance:      0.0,
		CreatedAt:    time.Now(),
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, "", err
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(ctx context.Context, req models.UserLogin) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret"
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.Hex(),
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(secret))
}

func (s *AuthService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *AuthService) UpdateProfile(ctx context.Context, userID primitive.ObjectID, firstName, lastName, email, phoneNumber string) error {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}

	if user.Email != email {
		existing, _ := s.userRepo.FindByEmail(ctx, email)
		if existing != nil {
			return errors.New("email already in use")
		}
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Email = email
	user.PhoneNumber = phoneNumber

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	newName := firstName + " " + lastName
	return s.reviewRepo.UpdateReviewerName(ctx, userID, newName)
}
