package services

import (
	"cinema-system/internal/models"
	"cinema-system/internal/repositories"
	"context"
	"errors"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingService struct {
	ticketRepo  *repositories.TicketRepository
	sessionRepo *repositories.SessionRepository
	userRepo    *repositories.UserRepository
	hallRepo    *repositories.HallRepository
	mu          sync.Mutex // For concurrent booking protection
}

func NewBookingService(
	ticketRepo *repositories.TicketRepository,
	sessionRepo *repositories.SessionRepository,
	userRepo *repositories.UserRepository,
	hallRepo *repositories.HallRepository,
) *BookingService {
	return &BookingService{
		ticketRepo:  ticketRepo,
		sessionRepo: sessionRepo,
		userRepo:    userRepo,
		hallRepo:    hallRepo,
	}
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

	ticket := &models.Ticket{
		SessionID:  sessionID,
		UserID:     userID,
		RowNumber:  row,
		SeatNumber: seat,
		Status:     models.TicketBooked,
		CreatedAt:  time.Now(),
	}

	errChan := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := s.ticketRepo.Create(ctx, ticket); err != nil {
			errChan <- err
		}
	}()

	go func() {
		defer wg.Done()
		newBalance := user.Balance - session.Price
		if err := s.userRepo.UpdateBalance(ctx, userID, newBalance); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	s.ticketRepo.UpdateStatus(ctx, ticket.ID, models.TicketPaid)

	return ticket, nil
}

func (s *BookingService) CancelTicket(ctx context.Context, ticketID, userID primitive.ObjectID) error {
	ticket, err := s.ticketRepo.FindByID(ctx, ticketID)
	if err != nil {
		return errors.New("ticket not found")
	}

	if ticket.UserID != userID {
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

	return s.ticketRepo.UpdateStatus(ctx, ticketID, models.TicketCancelled)
}

func (s *BookingService) GetUserTickets(ctx context.Context, userID primitive.ObjectID) ([]models.Ticket, error) {
	return s.ticketRepo.GetByUser(ctx, userID)
}

func (s *BookingService) GetSessionTickets(ctx context.Context, sessionID primitive.ObjectID) ([]models.Ticket, error) {
	return s.ticketRepo.GetBySession(ctx, sessionID)
}

func (s *BookingService) GetAllBookings(ctx context.Context) ([]models.Ticket, error) {
	return s.ticketRepo.GetAll(ctx)
}

// GetSessionBookedSeats returns row/seat pairs for occupied seats (for public seat map).
func (s *BookingService) GetSessionBookedSeats(ctx context.Context, sessionID primitive.ObjectID) ([]struct{ RowNumber int `json:"row_number"`; SeatNumber int `json:"seat_number"` }, error) {
	tickets, err := s.ticketRepo.GetBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	out := make([]struct{ RowNumber int `json:"row_number"`; SeatNumber int `json:"seat_number"` }, len(tickets))
	for i, t := range tickets {
		out[i].RowNumber = t.RowNumber
		out[i].SeatNumber = t.SeatNumber
	}
	return out, nil
}
