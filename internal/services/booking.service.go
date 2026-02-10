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
	movieRepo   *repositories.MovieRepository
	paymentRepo *repositories.PaymentRepository
	mu          sync.Mutex
}

func NewBookingService(
	ticketRepo *repositories.TicketRepository,
	sessionRepo *repositories.SessionRepository,
	userRepo *repositories.UserRepository,
	hallRepo *repositories.HallRepository,
	movieRepo *repositories.MovieRepository,
	paymentRepo *repositories.PaymentRepository,
) *BookingService {
	return &BookingService{
		ticketRepo:  ticketRepo,
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
		hallRepo:    hallRepo,
		movieRepo:   movieRepo,
		paymentRepo: paymentRepo,
	}
}

type SeatBookingRequest struct {
	RowNumber  int               `json:"row_number"`
	SeatNumber int               `json:"seat_number"`
	Type       models.TicketType `json:"type"`
}

func (s *BookingService) BookTickets(ctx context.Context, userID, sessionID primitive.ObjectID, seatRequests []SeatBookingRequest) ([]models.Ticket, error) {
	if len(seatRequests) == 0 {
		return nil, errors.New("no seats selected")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	session, err := s.sessionRepo.FindByID(ctx, sessionID)
	if err != nil {
		return nil, errors.New("session not found")
	}

	movie, err := s.movieRepo.FindByID(ctx, session.MovieID)
	if err != nil {
		return nil, errors.New("movie not found")
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	hall, err := s.hallRepo.FindByID(ctx, session.HallID)
	if err != nil {
		return nil, errors.New("hall not found")
	}

	var totalPrice float64
	var tickets []models.Ticket
	now := time.Now()

	for _, req := range seatRequests {
		if req.RowNumber < 1 || req.RowNumber > hall.TotalRows || req.SeatNumber < 1 || req.SeatNumber > hall.SeatsPerRow {
			return nil, fmt.Errorf("invalid seat position: row %d, seat %d", req.RowNumber, req.SeatNumber)
		}

		available, err := s.ticketRepo.CheckSeatAvailability(ctx, sessionID, req.RowNumber, req.SeatNumber)
		if err != nil {
			return nil, err
		}
		if !available {
			return nil, fmt.Errorf("seat already booked: row %d, seat %d", req.RowNumber, req.SeatNumber)
		}

		var priceMultiplier float64
		switch req.Type {
		case models.TicketAdult:
			priceMultiplier = 1.0
		case models.TicketStudent:
			priceMultiplier = 0.8
		case models.TicketPension:
			priceMultiplier = 0.7
		case models.TicketKid:
			if movie.AgeRating == "18+" {
				return nil, errors.New("kids tickets are not allowed for 18+ movies")
			}
			priceMultiplier = 0.5
		default:
			return nil, fmt.Errorf("invalid ticket type: %s", req.Type)
		}

		ticketPrice := session.Price * priceMultiplier
		totalPrice += ticketPrice

		tickets = append(tickets, models.Ticket{
			UserID:     userID,
			SessionID:  sessionID,
			RowNumber:  req.RowNumber,
			SeatNumber: req.SeatNumber,
			Type:       req.Type,
			Price:      ticketPrice,
			MovieTitle: movie.Name,
			Status:     models.TicketPaid,
			CreatedAt:  now,
		})
	}

	if user.Balance < totalPrice {
		return nil, fmt.Errorf("insufficient balance: need $%.2f, have $%.2f", totalPrice, user.Balance)
	}

	payment := &models.Payment{
		UserID:          userID,
		PaymentCardID:   primitive.NilObjectID,
		TransactionCode: s.generateTransactionCode(),
		Amount:          totalPrice,
		Status:          models.PaymentCompleted,
		CreatedAt:       now,
	}

	err = s.paymentRepo.Create(ctx, payment)
	if err != nil {
		return nil, errors.New("failed to create payment record")
	}

	newBalance := user.Balance - totalPrice
	err = s.userRepo.UpdateBalance(ctx, userID, newBalance)
	if err != nil {
		return nil, errors.New("failed to deduct balance")
	}

	for i := range tickets {
		tickets[i].PaymentID = payment.ID
		err = s.ticketRepo.Create(ctx, &tickets[i])
		if err != nil {
			return nil, fmt.Errorf("failed to create ticket for row %d, seat %d: %v", tickets[i].RowNumber, tickets[i].SeatNumber, err)
		}
	}

	return tickets, nil
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
	fmt.Printf("[DEBUG] GetUserTickets for userID: %s\n", userID.Hex())
	tickets, err := s.ticketRepo.GetByUserID(ctx, userID)
	if err != nil {
		fmt.Printf("[DEBUG] Error GetUserTickets: %v\n", err)
		return nil, err
	}
	fmt.Printf("[DEBUG] Found %d tickets\n", len(tickets))
	for i := range tickets {
		if tickets[i].MovieTitle == "" {
			session, _ := s.sessionRepo.FindByID(ctx, tickets[i].SessionID)
			if session != nil {
				movie, _ := s.movieRepo.FindByID(ctx, session.MovieID)
				if movie != nil {
					tickets[i].MovieTitle = movie.Name
				}
			}
		}
	}
	return tickets, nil
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
