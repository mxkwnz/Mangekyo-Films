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

	// Temporarily store genres to handle them after movie creation
	genreIDs := movie.Genres
	movie.Genres = []primitive.ObjectID{} // clear primarily if we want to rely solely on join table, but we can keep both for read optimization.
	// However, if we move to join table strict, we should probably only populate this on read.
	// For now, let's keep array for read performance but ALSO populate join table as requested.
	movie.Genres = genreIDs

	if err := s.movieRepo.Create(ctx, movie); err != nil {
		return err
	}

	for _, genreID := range genreIDs {
		if err := s.movieGenreService.AddGenreToMovie(ctx, movie.ID, genreID); err != nil {
			// In real app, rollback tx
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

	// Batch fetch genres to avoid N+1 queries
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
	// First update the movie document
	if err := s.movieRepo.Update(ctx, id, movie); err != nil {
		return err
	}

	// Update relationships
	// Simplest strategy: remove all existing and add new
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
