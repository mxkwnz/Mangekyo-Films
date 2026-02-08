package services

import (
	"cinema-system/internal/models"
	"cinema-system/internal/repositories"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieService struct {
	movieRepo *repositories.MovieRepository
	genreRepo *repositories.GenreRepository
}

func NewMovieService(movieRepo *repositories.MovieRepository, genreRepo *repositories.GenreRepository) *MovieService {
	return &MovieService{
		movieRepo: movieRepo,
		genreRepo: genreRepo,
	}
}

func (s *MovieService) CreateMovie(ctx context.Context, movie *models.Movie) error {
	movie.CreatedAt = time.Now()
	movie.Rating = 0.0
	return s.movieRepo.Create(ctx, movie)
}

func (s *MovieService) GetAllMovies(ctx context.Context) ([]models.Movie, error) {
	return s.movieRepo.GetAll(ctx)
}

func (s *MovieService) GetMovieByID(ctx context.Context, id primitive.ObjectID) (*models.Movie, error) {
	return s.movieRepo.FindByID(ctx, id)
}

func (s *MovieService) UpdateMovie(ctx context.Context, id primitive.ObjectID, movie *models.Movie) error {
	return s.movieRepo.Update(ctx, id, movie)
}

func (s *MovieService) DeleteMovie(ctx context.Context, id primitive.ObjectID) error {
	return s.movieRepo.Delete(ctx, id)
}
