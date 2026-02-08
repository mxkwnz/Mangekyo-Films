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

type ReviewService struct {
	reviewRepo *repositories.ReviewRepository
	movieRepo  *repositories.MovieRepository
	userRepo   *repositories.UserRepository
}

func NewReviewService(
	reviewRepo *repositories.ReviewRepository,
	movieRepo *repositories.MovieRepository,
	userRepo *repositories.UserRepository,
) *ReviewService {
	return &ReviewService{
		reviewRepo: reviewRepo,
		movieRepo:  movieRepo,
		userRepo:   userRepo,
	}
}

func (s *ReviewService) CreateReview(ctx context.Context, review *models.Review) error {
	exists, err := s.reviewRepo.CheckUserReview(ctx, review.UserID, review.MovieID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already reviewed this movie")
	}

	if review.Rating < 1 || review.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	_, err = s.movieRepo.FindByID(ctx, review.MovieID)
	if err != nil {
		return errors.New("movie not found")
	}

	_, err = s.userRepo.FindByID(ctx, review.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	review.CreatedAt = time.Now()

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return err
	}

	go s.updateMovieRating(context.Background(), review.MovieID)

	return nil
}

func (s *ReviewService) updateMovieRating(ctx context.Context, movieID primitive.ObjectID) {
	avgRating, err := s.reviewRepo.GetAverageRating(ctx, movieID)
	if err != nil {
		return
	}

	s.movieRepo.UpdateRating(ctx, movieID, avgRating)
}

func (s *ReviewService) GetMovieReviews(ctx context.Context, movieID primitive.ObjectID) ([]models.Review, error) {
	return s.reviewRepo.GetByMovie(ctx, movieID)
}

func (s *ReviewService) DeleteReview(ctx context.Context, reviewID, userID primitive.ObjectID, isAdmin bool) error {
	review, err := s.reviewRepo.FindByID(ctx, reviewID)
	if err != nil {
		return errors.New("review not found")
	}

	if !isAdmin && review.UserID != userID {
		return errors.New("unauthorized to delete this review")
	}

	if err := s.reviewRepo.Delete(ctx, reviewID); err != nil {
		return err
	}

	go s.updateMovieRating(context.Background(), review.MovieID)

	return nil
}

func (s *ReviewService) CalculateMovieRating(ctx context.Context, movieID primitive.ObjectID) (float64, error) {
	return s.reviewRepo.GetAverageRating(ctx, movieID)
}

func (s *ReviewService) BatchUpdateRatings(ctx context.Context, movieIDs []primitive.ObjectID) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(movieIDs))

	for _, movieID := range movieIDs {
		wg.Add(1)
		go func(id primitive.ObjectID) {
			defer wg.Done()
			avgRating, err := s.reviewRepo.GetAverageRating(ctx, id)
			if err != nil {
				errChan <- err
				return
			}
			if err := s.movieRepo.UpdateRating(ctx, id, avgRating); err != nil {
				errChan <- err
			}
		}(movieID)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
