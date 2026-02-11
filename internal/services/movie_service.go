package services

import (
	"cinema-system/internal/models"
	"cinema-system/internal/repositories"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MovieService struct {
	movieRepo         *repositories.MovieRepository
	genreRepo         *repositories.GenreRepository
	movieGenreService *MovieGenreService
}

func NewMovieService(
	movieRepo *repositories.MovieRepository,
	genreRepo *repositories.GenreRepository,
	movieGenreService *MovieGenreService,
) *MovieService {
	return &MovieService{
		movieRepo:         movieRepo,
		genreRepo:         genreRepo,
		movieGenreService: movieGenreService,
	}
}

func (s *MovieService) CreateMovie(ctx context.Context, movie *models.Movie) error {
	movie.CreatedAt = time.Now()
	movie.Rating = 0.0

	genreIDs := movie.Genres
	movie.Genres = genreIDs

	if err := s.movieRepo.Create(ctx, movie); err != nil {
		return err
	}

	for _, genreID := range genreIDs {
		if err := s.movieGenreService.AddGenreToMovie(ctx, movie.ID, genreID); err != nil {
			return err
		}
	}
	return nil
}

func (s *MovieService) GetAllMovies(ctx context.Context) ([]models.Movie, error) {
	movies, err := s.movieRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	genreIDMap := make(map[primitive.ObjectID]bool)
	for _, m := range movies {
		for _, gid := range m.Genres {
			genreIDMap[gid] = true
		}
	}

	uniqueGenreIDs := make([]primitive.ObjectID, 0, len(genreIDMap))
	for gid := range genreIDMap {
		uniqueGenreIDs = append(uniqueGenreIDs, gid)
	}

	genres, err := s.genreRepo.FindByIDs(ctx, uniqueGenreIDs)
	if err != nil {
		return nil, err
	}

	genreMap := make(map[primitive.ObjectID]string)
	for _, g := range genres {
		genreMap[g.ID] = g.Name
	}

	for i := range movies {
		movie := &movies[i]
		names := make([]string, 0, len(movie.Genres))
		for _, gid := range movie.Genres {
			if name, ok := genreMap[gid]; ok {
				names = append(names, name)
			}
		}
		movie.GenreNames = names
	}

	return movies, nil
}

func (s *MovieService) GetMovieByID(ctx context.Context, id primitive.ObjectID) (*models.Movie, error) {
	movie, err := s.movieRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.populateGenreNames(ctx, movie); err != nil {
		return nil, err
	}
	return movie, nil
}

func (s *MovieService) populateGenreNames(ctx context.Context, movie *models.Movie) error {
	if len(movie.Genres) == 0 {
		movie.GenreNames = []string{}
		return nil
	}
	genres, err := s.genreRepo.FindByIDs(ctx, movie.Genres)
	if err != nil {
		return err
	}
	names := make([]string, len(genres))
	for i, g := range genres {
		names[i] = g.Name
	}
	movie.GenreNames = names
	return nil
}

func (s *MovieService) UpdateMovie(ctx context.Context, id primitive.ObjectID, movie *models.Movie) error {
	if err := s.movieRepo.Update(ctx, id, movie); err != nil {
		return err
	}

	if err := s.movieGenreService.RemoveGenresFromMovie(ctx, id); err != nil {
		return err
	}

	for _, genreID := range movie.Genres {
		if err := s.movieGenreService.AddGenreToMovie(ctx, id, genreID); err != nil {
			return err
		}
	}

	return nil
}

func (s *MovieService) DeleteMovie(ctx context.Context, id primitive.ObjectID) error {
	if err := s.movieGenreService.RemoveGenresFromMovie(ctx, id); err != nil {
		return err
	}
	return s.movieRepo.Delete(ctx, id)
}
