package routes

import (
	"net/http"

	"cinema-system/internal/handlers"
)

func RegisterRoutes(
	movieHandler *handlers.MovieHandler,
	ticketHandler *handlers.TicketHandler,
	sessionHandler *handlers.SessionHandler,
) {
	http.HandleFunc("/movies", movieHandler.HandleMovies)
	http.HandleFunc("/tickets", ticketHandler.HandleTickets)
	http.HandleFunc("/sessions", sessionHandler.HandleSessions)
}
