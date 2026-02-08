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
	movie, err := s.movieRepo.FindByID(ctx, session.MovieID)
	if err != nil {
		return errors.New("movie not found")
	}

	// Validate hall exists
	_, err = s.hallRepo.FindByID(ctx, session.HallID)
	if err != nil {
		return errors.New("hall not found")
	}

	session.EndTime = session.StartTime.Add(time.Duration(movie.Duration) * time.Minute)

	return s.sessionRepo.Create(ctx, session)
}

func (s *SessionService) GetSessionsByMovie(ctx context.Context, movieID primitive.ObjectID) ([]models.Session, error) {
	return s.sessionRepo.GetByMovie(ctx, movieID)
}

func (s *SessionService) GetUpcomingSessions(ctx context.Context) ([]models.Session, error) {
	return s.sessionRepo.GetUpcoming(ctx)
}

func (s *SessionService) GetSessionByID(ctx context.Context, id primitive.ObjectID) (*models.Session, error) {
	return s.sessionRepo.FindByID(ctx, id)
}

func (s *SessionService) DeleteSession(ctx context.Context, id primitive.ObjectID) error {
	return s.sessionRepo.Delete(ctx, id)
}
