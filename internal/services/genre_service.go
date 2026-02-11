package services

import (
	"cinema-system/internal/models"
	"cinema-system/internal/repositories"
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GenreService struct {
	genreRepo *repositories.GenreRepository
}

func NewGenreService(genreRepo *repositories.GenreRepository) *GenreService {
	return &GenreService{genreRepo: genreRepo}
}

func (s *GenreService) CreateGenre(ctx context.Context, genre *models.Genre) error {
	return s.genreRepo.Create(ctx, genre)
}

func (s *GenreService) GetAllGenres(ctx context.Context) ([]models.Genre, error) {
	return s.genreRepo.GetAll(ctx)
}

func (s *GenreService) GetGenreByID(ctx context.Context, id primitive.ObjectID) (*models.Genre, error) {
	return s.genreRepo.FindByID(ctx, id)
}

func (s *GenreService) UpdateGenre(ctx context.Context, id primitive.ObjectID, genre *models.Genre) error {
	return s.genreRepo.Update(ctx, id, genre)
}

func (s *GenreService) DeleteGenre(ctx context.Context, id primitive.ObjectID) error {
	return s.genreRepo.Delete(ctx, id)
}
