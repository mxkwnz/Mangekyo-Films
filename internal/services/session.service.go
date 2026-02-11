package services

import (
	"cinema-system/internal/models"
	"cinema-system/internal/repositories"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SessionService struct {
	sessionRepo *repositories.SessionRepository
	hallRepo    *repositories.HallRepository
	movieRepo   *repositories.MovieRepository
}

func NewSessionService(
	sessionRepo *repositories.SessionRepository,
	hallRepo *repositories.HallRepository,
	movieRepo *repositories.MovieRepository,
) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
		hallRepo:    hallRepo,
		movieRepo:   movieRepo,
	}
}

func (s *SessionService) CreateSession(ctx context.Context, session *models.Session) error {
	if session.MovieID.IsZero() {
		return errors.New("Movie ID is required")
	}
	if session.HallID.IsZero() {
		return errors.New("Hall ID is required")
	}

	movie, err := s.movieRepo.FindByID(ctx, session.MovieID)
	if err != nil {
		return errors.New("Movie does not exist")
	}

	_, err = s.hallRepo.FindByID(ctx, session.HallID)
	if err != nil {
		return errors.New("Hall does not exist")
	}

	// Allow 5 minutes buffer for minor clock drift
	if session.StartTime.Before(time.Now().Add(-5 * time.Minute)) {
		return errors.New("Cannot schedule sessions in the past")
	}

	session.EndTime = session.StartTime.Add(time.Duration(movie.Duration) * time.Minute)

	overlapping, err := s.sessionRepo.GetOverlappingByHall(ctx, session.HallID, session.StartTime, session.EndTime)
	if err != nil {
		return err
	}
	if len(overlapping) > 0 {
		return errors.New("The selected hall is already occupied during this time period")
	}

	return s.sessionRepo.Create(ctx, session)
}

func (s *SessionService) GetSessionsByMovie(ctx context.Context, movieID primitive.ObjectID) ([]models.Session, error) {
	return s.sessionRepo.GetByMovie(ctx, movieID)
}

func (s *SessionService) GetUpcomingSessions(ctx context.Context) ([]models.Session, error) {
	return s.sessionRepo.GetUpcoming(ctx)
}

func (s *SessionService) GetUpcomingMovieIDs(ctx context.Context) ([]primitive.ObjectID, error) {
	return s.sessionRepo.GetUpcomingMovieIDs(ctx)
}

func (s *SessionService) GetSessionByID(ctx context.Context, id primitive.ObjectID) (*models.Session, error) {
	return s.sessionRepo.FindByID(ctx, id)
}

func (s *SessionService) UpdateSession(ctx context.Context, id primitive.ObjectID, session *models.Session) error {
	if session.MovieID.IsZero() {
		return errors.New("Movie ID is required")
	}
	if session.HallID.IsZero() {
		return errors.New("Hall ID is required")
	}

	movie, err := s.movieRepo.FindByID(ctx, session.MovieID)
	if err != nil {
		return errors.New("Movie does not exist")
	}

	_, err = s.hallRepo.FindByID(ctx, session.HallID)
	if err != nil {
		return errors.New("Hall does not exist")
	}

	// We allow updating sessions in the past if they were already there, but let's keep some sanity check.
	// Usually, you only update future sessions.
	if session.StartTime.Before(time.Now().Add(-5 * time.Minute)) {
		// If it's a minor update and start time hasn't changed much, maybe it's fine.
		// For simplicity, let's keep the past check.
		return errors.New("Cannot schedule sessions in the past")
	}

	session.EndTime = session.StartTime.Add(time.Duration(movie.Duration) * time.Minute)

	overlapping, err := s.sessionRepo.GetOverlappingByHall(ctx, session.HallID, session.StartTime, session.EndTime)
	if err != nil {
		return err
	}

	// Filter out the current session from overlap check
	for _, o := range overlapping {
		if o.ID != id {
			return errors.New("The selected hall is already occupied during this time period")
		}
	}

	return s.sessionRepo.Update(ctx, id, session)
}

func (s *SessionService) DeleteSession(ctx context.Context, id primitive.ObjectID) error {
	return s.sessionRepo.Delete(ctx, id)
}
