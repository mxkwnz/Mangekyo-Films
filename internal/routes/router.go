package routes

import (
	"net/http"

	"cinema-system/internal/handlers"
)

func RegisterRoutes(
	movieHandler *handlers.MovieHandler,
	ticketHandler *handlers.TicketHandler,
) {
	http.HandleFunc("/movies", movieHandler.HandleMovies)
	http.HandleFunc("/tickets", ticketHandler.HandleTickets)
}
