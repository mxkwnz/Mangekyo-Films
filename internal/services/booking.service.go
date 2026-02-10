package services

import (
	"cinema-system/internal/models"
	"cinema-system/internal/repositories"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingService struct {
	ticketRepo  *repositories.TicketRepository
	sessionRepo *repositories.SessionRepository
	userRepo    *repositories.UserRepository
	hallRepo    *repositories.HallRepository
	paymentRepo *repositories.PaymentRepository
	mu          sync.Mutex // For concurrent booking protection
}

type SeatBookingRequest struct {
	RowNumber  int
	SeatNumber int
	Type       models.TicketType
}

func NewBookingService(
	ticketRepo *repositories.TicketRepository,
	sessionRepo *repositories.SessionRepository,
	userRepo *repositories.UserRepository,
	hallRepo *repositories.HallRepository,
	paymentRepo *repositories.PaymentRepository,
) *BookingService {
	return &BookingService{
		ticketRepo:  ticketRepo,
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
		hallRepo:    hallRepo,
		paymentRepo: paymentRepo,
	}
}

func (s *BookingService) BookTickets(ctx context.Context, userID, sessionID primitive.ObjectID, seats []SeatBookingRequest) ([]*models.Ticket, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 1. Validate Session
	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, errors.New("session not found")
	}

	// 2. Validate User & Balance
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	totalPrice := float64(len(seats)) * session.Price
	if user.Balance < totalPrice {
		return nil, errors.New("insufficient balance")
	}

	// 3. Validate Hall & Seats
	hall, err := s.hallRepo.FindByID(ctx, session.HallID)
	if err != nil {
		return nil, errors.New("hall not found")
	}

	for _, seat := range seats {
		if seat.RowNumber < 1 || seat.RowNumber > hall.TotalRows || seat.SeatNumber < 1 || seat.SeatNumber > hall.SeatsPerRow {
			return nil, fmt.Errorf("invalid seat position: row %d, seat %d", seat.RowNumber, seat.SeatNumber)
		}

		available, err := s.ticketRepo.CheckSeatAvailability(ctx, sessionID, seat.RowNumber, seat.SeatNumber)
		if err != nil {
			return nil, err
		}
		if !available {
			return nil, fmt.Errorf("seat row %d number %d already booked", seat.RowNumber, seat.SeatNumber)
		}
	}

	// 4. Create Payment
	payment := &models.Payment{
		UserID:          userID,
		PaymentCardID:   primitive.NilObjectID,
		TransactionCode: s.generateTransactionCode(),
		Amount:          totalPrice,
		Status:          models.PaymentPending,
		CreatedAt:       time.Now(),
	}

	err = s.paymentRepo.Create(ctx, payment)
	if err != nil {
		return nil, errors.New("failed to create payment record")
	}

	// 5. Deduct Balance
	newBalance := user.Balance - totalPrice
	err = s.userRepo.UpdateBalance(ctx, userID, newBalance)
	if err != nil {
		return nil, errors.New("failed to deduct balance")
	}

	// 6. Confirm Payment
	err = s.paymentRepo.UpdateStatus(ctx, payment.ID, models.PaymentCompleted)
	if err != nil {
		return nil, err
	}

	// 7. Create Tickets
	var tickets []*models.Ticket
	for _, seat := range seats {
		ticket := &models.Ticket{
			UserID:     userID,
			SessionID:  sessionID,
			PaymentID:  payment.ID,
			RowNumber:  seat.RowNumber,
			SeatNumber: seat.SeatNumber,
			Type:       seat.Type,
			Price:      session.Price,
			MovieTitle: "", // populated if needed, or left empty
			Status:     models.TicketPaid,
			CreatedAt:  time.Now(),
		}

		// Attempt to fetch movie title if we can, but it's okay if not for now
		// In a real app we might fetch the movie from session.MovieID

		err = s.ticketRepo.Create(ctx, ticket)
		if err != nil {
			// In a real system we would need to rollback payment here
			return nil, fmt.Errorf("failed to create ticket for row %d seat %d: %w", seat.RowNumber, seat.SeatNumber, err)
		}
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

func (s *BookingService) BookTicket(ctx context.Context, userID, sessionID primitive.ObjectID, row, seat int) (*models.Ticket, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	available, err := s.ticketRepo.CheckSeatAvailability(ctx, sessionID, row, seat)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, errors.New("seat already booked")
	}

	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, errors.New("session not found")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Balance < session.Price {
		return nil, errors.New("insufficient balance")
	}

	hall, err := s.hallRepo.FindByID(ctx, session.HallID)
	if err != nil {
		return nil, errors.New("hall not found")
	}

	if row < 1 || row > hall.TotalRows || seat < 1 || seat > hall.SeatsPerRow {
		return nil, errors.New("invalid seat position")
	}

	payment := &models.Payment{
		UserID:          userID,
		PaymentCardID:   primitive.NilObjectID,
		TransactionCode: s.generateTransactionCode(),
		Amount:          session.Price,
		Status:          models.PaymentPending,
		CreatedAt:       time.Now(),
	}

	err = s.paymentRepo.Create(ctx, payment)
	if err != nil {
		return nil, errors.New("failed to create payment record")
	}

	newBalance := user.Balance - session.Price
	err = s.userRepo.UpdateBalance(ctx, userID, newBalance)
	if err != nil {
		return nil, errors.New("failed to deduct balance")
	}

	err = s.paymentRepo.UpdateStatus(ctx, payment.ID, models.PaymentCompleted)
	if err != nil {
		return nil, err
	}

	ticket := &models.Ticket{
		SessionID:  sessionID,
		PaymentID:  payment.ID,
		RowNumber:  row,
		SeatNumber: seat,
		Status:     models.TicketPaid,
		CreatedAt:  time.Now(),
	}

	err = s.ticketRepo.Create(ctx, ticket)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *BookingService) generateTransactionCode() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("TXN-%d-%s", timestamp, primitive.NewObjectID().Hex()[:8])
}

func (s *BookingService) CancelTicket(ctx context.Context, ticketID, userID primitive.ObjectID) error {
	ticket, err := s.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return errors.New("ticket not found")
	}

	payment, err := s.paymentRepo.FindByID(ctx, ticket.PaymentID)
	if err != nil {
		return errors.New("payment not found")
	}

	if payment.UserID != userID {
		return errors.New("unauthorized to cancel this ticket")
	}

	if ticket.Status == models.TicketCancelled {
		return errors.New("ticket already cancelled")
	}

	session, err := s.sessionRepo.FindByID(ctx, ticket.SessionID)
	if err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	newBalance := user.Balance + session.Price
	if err := s.userRepo.UpdateBalance(ctx, userID, newBalance); err != nil {
		return err
	}

	if err := s.paymentRepo.UpdateStatus(ctx, payment.ID, models.PaymentRefunded); err != nil {
		return err
	}

	return s.ticketRepo.UpdateStatus(ctx, ticketID, models.TicketCancelled)
}

func (s *BookingService) GetUserTickets(ctx context.Context, userID primitive.ObjectID) ([]models.Ticket, error) {
	payments, err := s.paymentRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(payments) == 0 {
		return []models.Ticket{}, nil
	}

	paymentIDs := make([]primitive.ObjectID, len(payments))
	for i, payment := range payments {
		paymentIDs[i] = payment.ID
	}

	return s.ticketRepo.GetByPaymentIDs(ctx, paymentIDs)
}

func (s *BookingService) GetSessionTickets(ctx context.Context, sessionID primitive.ObjectID) ([]models.Ticket, error) {
	return s.ticketRepo.GetBySession(ctx, sessionID)
}

func (s *BookingService) GetAllBookings(ctx context.Context) ([]models.Ticket, error) {
	return s.ticketRepo.GetAll(ctx)
}

func (s *BookingService) GetSessionBookedSeats(ctx context.Context, sessionID primitive.ObjectID) ([]struct {
	RowNumber  int `json:"row_number"`
	SeatNumber int `json:"seat_number"`
}, error) {
	tickets, err := s.ticketRepo.GetBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	out := make([]struct {
		RowNumber  int `json:"row_number"`
		SeatNumber int `json:"seat_number"`
	}, len(tickets))
	for i, t := range tickets {
		out[i].RowNumber = t.RowNumber
		out[i].SeatNumber = t.SeatNumber
	}
	return out, nil
}
