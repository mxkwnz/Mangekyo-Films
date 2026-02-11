package services

import (
	"cinema-system/internal/models"
	"cinema-system/internal/repositories"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type PaymentCardService struct {
	cardRepo *repositories.PaymentCardRepository
	userRepo *repositories.UserRepository
}

func NewPaymentCardService(cardRepo *repositories.PaymentCardRepository, userRepo *repositories.UserRepository) *PaymentCardService {
	return &PaymentCardService{
		cardRepo: cardRepo,
		userRepo: userRepo,
	}
}

func (s *PaymentCardService) CreateCard(ctx context.Context, userID primitive.ObjectID, req models.PaymentCardCreate) (*models.PaymentCard, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// CVV Hashing
	cvvHash, err := bcrypt.GenerateFromPassword([]byte(req.CVV), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to process card security")
	}

	card := &models.PaymentCard{
		UserID:         userID,
		CardHolderName: req.CardHolderName,
		CardNumber:     req.CardNumber,
		ExpiryDate:     req.ExpiryDate,
		CVVHash:        string(cvvHash),
		CreatedAt:      time.Now(),
	}

	err = s.cardRepo.Create(ctx, card)
	if err != nil {
		return nil, err
	}

	return card, nil
}

func (s *PaymentCardService) GetUserCards(ctx context.Context, userID primitive.ObjectID) ([]models.PaymentCard, error) {
	cards, err := s.cardRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (s *PaymentCardService) GetCardByID(ctx context.Context, cardID primitive.ObjectID, userID primitive.ObjectID) (*models.PaymentCard, error) {
	card, err := s.cardRepo.FindByID(ctx, cardID)
	if err != nil {
		return nil, errors.New("card not found")
	}

	if card.UserID != userID {
		return nil, errors.New("unauthorized access to card")
	}

	return card, nil
}

func (s *PaymentCardService) DeleteCard(ctx context.Context, cardID primitive.ObjectID, userID primitive.ObjectID) error {
	card, err := s.cardRepo.FindByID(ctx, cardID)
	if err != nil {
		return errors.New("card not found")
	}

	if card.UserID != userID {
		return errors.New("unauthorized access to card")
	}

	return s.cardRepo.Delete(ctx, cardID)
}
