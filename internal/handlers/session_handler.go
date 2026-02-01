package handlers

import (
	"cinema-system/internal/models"
	"cinema-system/internal/repositories"
	"encoding/json"
	"net/http"
)

type SessionHandler struct {
	repo *repositories.SessionRepo
}

func NewSessionHandler(repo *repositories.SessionRepo) *SessionHandler {
	return &SessionHandler{repo: repo}
}

func (h *SessionHandler) HandleSessions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		json.NewEncoder(w).Encode(h.repo.GetAll())
		return
	}

	if r.Method == "POST" {
		var session models.Session
		json.NewDecoder(r.Body).Decode(&session)

		created := h.repo.Create(session)
		json.NewEncoder(w).Encode(created)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
