package services

import (
	"cinema-system/internal/models"
	"cinema-system/internal/repositories"
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentService struct {
	paymentRepo *repositories.PaymentRepository
	cardRepo    *repositories.PaymentCardRepository
	userRepo    *repositories.UserRepository
}

func NewPaymentService(
	paymentRepo *repositories.PaymentRepository,
	cardRepo *repositories.PaymentCardRepository,
	userRepo *repositories.UserRepository,
) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		cardRepo:    cardRepo,
		userRepo:    userRepo,
	}
}

func (s *PaymentService) generateTransactionCode() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("TXN-%d-%s", timestamp, primitive.NewObjectID().Hex()[:8])
}

func (s *PaymentService) CreatePayment(ctx context.Context, userID primitive.ObjectID, req models.PaymentCreate) (*models.Payment, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	card, err := s.cardRepo.FindByID(ctx, req.PaymentCardID)
	if err != nil {
		return nil, errors.New("payment card not found")
	}

	if card.UserID != userID {
		return nil, errors.New("unauthorized access to payment card")
	}

	if user.Balance < req.Amount {
		return nil, errors.New("insufficient balance")
	}

	payment := &models.Payment{
		UserID:          userID,
		PaymentCardID:   req.PaymentCardID,
		TransactionCode: s.generateTransactionCode(),
		Amount:          req.Amount,
		Status:          models.PaymentPending,
		CreatedAt:       time.Now(),
	}

	err = s.paymentRepo.Create(ctx, payment)
	if err != nil {
		return nil, err
	}

	newBalance := user.Balance - req.Amount
	err = s.userRepo.UpdateBalance(ctx, userID, newBalance)
	if err != nil {
		// Rollback would be needed here in production
		return nil, errors.New("failed to process payment")
	}

	err = s.paymentRepo.UpdateStatus(ctx, payment.ID, models.PaymentCompleted)
	if err != nil {
		return nil, err
	}
	payment.Status = models.PaymentCompleted

	return payment, nil
}

func (s *PaymentService) GetPaymentByID(ctx context.Context, paymentID primitive.ObjectID, userID primitive.ObjectID) (*models.Payment, error) {
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return nil, errors.New("payment not found")
	}

	if payment.UserID != userID {
		return nil, errors.New("unauthorized access to payment")
	}

	return payment, nil
}

func (s *PaymentService) GetUserPayments(ctx context.Context, userID primitive.ObjectID) ([]models.Payment, error) {
	payments, err := s.paymentRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return payments, nil
}

func (s *PaymentService) GetAllPayments(ctx context.Context) ([]models.Payment, error) {
	return s.paymentRepo.GetAll(ctx)
}

func (s *PaymentService) RefundPayment(ctx context.Context, paymentID primitive.ObjectID, userID primitive.ObjectID) error {
	payment, err := s.paymentRepo.FindByID(ctx, paymentID)
	if err != nil {
		return errors.New("payment not found")
	}

	if payment.UserID != userID {
		return errors.New("unauthorized access to payment")
	}

	if payment.Status != models.PaymentCompleted {
		return errors.New("only completed payments can be refunded")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	newBalance := user.Balance + payment.Amount
	err = s.userRepo.UpdateBalance(ctx, userID, newBalance)
	if err != nil {
		return errors.New("failed to process refund")
	}

	err = s.paymentRepo.UpdateStatus(ctx, paymentID, models.PaymentRefunded)
	if err != nil {
		return err
	}

	return nil
}

// TopUpBalance adds money to user's balance using a payment card
func (s *PaymentService) TopUpBalance(ctx context.Context, userID primitive.ObjectID, req models.PaymentCreate) (*models.Payment, error) {
	// Validate payment request
	if err := req.Validate(); err != nil {
		return nil, err
	}

	// Verify user exists
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Verify payment card exists and belongs to user
	card, err := s.cardRepo.FindByID(ctx, req.PaymentCardID)
	if err != nil {
		return nil, errors.New("payment card not found")
	}

	if card.UserID != userID {
		return nil, errors.New("unauthorized access to payment card")
	}

	// Create payment record for top-up
	payment := &models.Payment{
		UserID:          userID,
		PaymentCardID:   req.PaymentCardID,
		TransactionCode: s.generateTransactionCode(),
		Amount:          req.Amount,
		Status:          models.PaymentPending,
		CreatedAt:       time.Now(),
	}

	err = s.paymentRepo.Create(ctx, payment)
	if err != nil {
		return nil, err
	}

	// Add amount to user balance
	newBalance := user.Balance + req.Amount
	err = s.userRepo.UpdateBalance(ctx, userID, newBalance)
	if err != nil {
		return nil, errors.New("failed to process top-up")
	}

	// Update payment status to completed
	err = s.paymentRepo.UpdateStatus(ctx, payment.ID, models.PaymentCompleted)
	if err != nil {
		return nil, err
	}
	payment.Status = models.PaymentCompleted

	return payment, nil
}
