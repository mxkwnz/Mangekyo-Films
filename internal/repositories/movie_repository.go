package repositories

import (
	"cinema-system/internal/models"
	"sync"
	"time"
)

type MovieRepo struct {
	mu     sync.RWMutex
	data   map[int]models.Movie
	nextID int
}

func NewMovieRepo() *MovieRepo {
	return &MovieRepo{
		data:   make(map[int]models.Movie),
		nextID: 1,
	}
}

func (r *MovieRepo) Create(movie models.Movie) models.Movie {
	r.mu.Lock()
	defer r.mu.Unlock()

	movie.ID = r.nextID
	movie.CreatedAt = time.Now()
	r.nextID++

	r.data[movie.ID] = movie
	return movie
}

func (r *MovieRepo) GetAll() []models.Movie {
	r.mu.RLock()
	defer r.mu.RUnlock()

	list := []models.Movie{}
	for _, m := range r.data {
		list = append(list, m)
	}
	return list
}
