package handlers

import (
	"context"
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

func NewTicketHandler(
	repo *repositories.TicketRepository,
	worker *services.BookingWorker,
) *TicketHandler {
	return &TicketHandler{repo: repo, worker: worker}
}

func (h *TicketHandler) HandleTickets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := context.Background()

	switch r.Method {

	case http.MethodGet:
		tickets, err := h.repo.GetAll(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(tickets)

	case http.MethodPost:
		var ticket models.Ticket
		if err := json.NewDecoder(r.Body).Decode(&ticket); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.repo.Create(ctx, &ticket); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.worker.Queue <- "Ticket booked successfully"
		json.NewEncoder(w).Encode(ticket)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
