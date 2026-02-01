package handlers

import (
	"cinema-system/internal/models"
	"cinema-system/internal/repositories"
	"encoding/json"
	"net/http"
)

type MovieHandler struct {
	repo *repositories.MovieRepo
}

func NewMovieHandler(repo *repositories.MovieRepo) *MovieHandler {
	return &MovieHandler{repo: repo}
}

func (h *MovieHandler) HandleMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		json.NewEncoder(w).Encode(h.repo.GetAll())
		return
	}

	if r.Method == "POST" {
		var movie models.Movie
		json.NewDecoder(r.Body).Decode(&movie)

		created := h.repo.Create(movie)
		json.NewEncoder(w).Encode(created)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
