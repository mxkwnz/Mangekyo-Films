package services

import (
	"cinema-system/internal/models"
	"cinema-system/internal/repositories"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieGenreService struct {
	movieGenreRepo *repositories.MovieGenreRepository
}

func NewMovieGenreService(movieGenreRepo *repositories.MovieGenreRepository) *MovieGenreService {
	return &MovieGenreService{movieGenreRepo: movieGenreRepo}
}

func (s *MovieGenreService) AddGenreToMovie(ctx context.Context, movieID, genreID primitive.ObjectID) error {
	movieGenre := &models.MovieGenre{
		MovieID: movieID,
		GenreID: genreID,
	}
	return s.movieGenreRepo.Create(ctx, movieGenre)
}

func (s *MovieGenreService) GetGenresByMovieID(ctx context.Context, movieID primitive.ObjectID) ([]models.MovieGenre, error) {
	return s.movieGenreRepo.GetGenresByMovieID(ctx, movieID)
}

func (s *MovieGenreService) RemoveGenresFromMovie(ctx context.Context, movieID primitive.ObjectID) error {
	return s.movieGenreRepo.DeleteByMovieID(ctx, movieID)
}
