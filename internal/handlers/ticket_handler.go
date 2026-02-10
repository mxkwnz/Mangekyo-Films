package handlers

import (
	"encoding/json"
	"net/http"

	"cinema-system/internal/models"
	"cinema-system/internal/repositories"
	"cinema-system/internal/services"
)

type TicketHandler struct {
	repo   *repositories.TicketRepository
	worker *services.BookingWorker
}

func NewTicketHandler(repo *repositories.TicketRepository, worker *services.BookingWorker) *TicketHandler {
	return &TicketHandler{repo: repo, worker: worker}
}

func (h *TicketHandler) HandleTickets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	if r.Method == "GET" {
		tickets, err := h.repo.GetAll(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}
		json.NewEncoder(w).Encode(tickets)
		return
	}

	if r.Method == "POST" {
		var ticket models.Ticket
		if err := json.NewDecoder(r.Body).Decode(&ticket); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
			return
		}

		if err := h.repo.Create(ctx, &ticket); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		h.worker.Queue <- "Ticket booked successfully"

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(ticket)
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}
