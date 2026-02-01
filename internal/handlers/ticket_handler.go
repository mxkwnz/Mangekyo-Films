package handlers

import (
	"encoding/json"
	"net/http"

	"cinema-system/internal/models"
	"cinema-system/internal/repositories"
	"cinema-system/internal/services"
)

type TicketHandler struct {
	repo   *repositories.TicketRepo
	worker *services.BookingWorker
}

func NewTicketHandler(repo *repositories.TicketRepo, worker *services.BookingWorker) *TicketHandler {
	return &TicketHandler{repo: repo, worker: worker}
}

func (h *TicketHandler) HandleTickets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		json.NewEncoder(w).Encode(h.repo.GetAll())
		return
	}

	if r.Method == "POST" {
		var ticket models.Ticket
		json.NewDecoder(r.Body).Decode(&ticket)

		booked := h.repo.Book(ticket)

		h.worker.Queue <- "Ticket booked successfully"

		json.NewEncoder(w).Encode(booked)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
