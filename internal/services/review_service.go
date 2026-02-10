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

	if review.Rating < 0 || review.Rating > 10 {
		return errors.New("rating must be between 0 and 10")
	}

	movie, err := s.movieRepo.FindByID(ctx, review.MovieID)
	if err != nil {
		return errors.New("movie not found")
	}

	user, err := s.userRepo.FindByID(ctx, review.UserID)
	if err != nil {
		return errors.New("user not found")
	}

	review.UserName = user.FirstName + " " + user.LastName
	review.MovieTitle = movie.Name
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
	reviews, err := s.reviewRepo.GetByMovie(ctx, movieID)
	if err != nil {
		return nil, err
	}

	for i := range reviews {
		if reviews[i].UserName == "" {
			u, _ := s.userRepo.FindByID(ctx, reviews[i].UserID)
			if u != nil {
				reviews[i].UserName = u.FirstName + " " + u.LastName
			}
		}
	}
	return reviews, nil
}

func (s *ReviewService) GetMyReviews(ctx context.Context, userID primitive.ObjectID) ([]models.Review, error) {
	fmt.Printf("[DEBUG] GetMyReviews for userID: %s\n", userID.Hex())
	reviews, err := s.reviewRepo.GetByUser(ctx, userID)
	if err != nil {
		fmt.Printf("[DEBUG] Error GetMyReviews: %v\n", err)
		return nil, err
	}
	fmt.Printf("[DEBUG] Found %d reviews\n", len(reviews))

	for i := range reviews {
		if reviews[i].MovieTitle == "" {
			m, _ := s.movieRepo.FindByID(ctx, reviews[i].MovieID)
			if m != nil {
				reviews[i].MovieTitle = m.Name
			}
		}
	}
	return reviews, nil
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

func (s *ReviewService) UpdateReview(ctx context.Context, reviewID, userID primitive.ObjectID, rating int, comment string) error {
	review, err := s.reviewRepo.FindByID(ctx, reviewID)
	if err != nil {
		return errors.New("review not found")
	}

	if review.UserID != userID {
		return errors.New("unauthorized to update this review")
	}

	if rating < 0 || rating > 10 {
		return errors.New("rating must be between 0 and 10")
	}

	review.Rating = rating
	review.Comment = comment
	review.CreatedAt = time.Now()

	if err := s.reviewRepo.Update(ctx, review); err != nil {
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
