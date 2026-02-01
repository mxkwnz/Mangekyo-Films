package main

import (
	"cinema-system/internal/handlers"
	"cinema-system/internal/repositories"
	"cinema-system/internal/routes"
	"cinema-system/internal/services"
	"log"
	"net/http"
)

func main() {
	movieRepo := repositories.NewMovieRepo()
	sessionRepo := repositories.NewSessionRepo()
	ticketRepo := repositories.NewTicketRepo()

	worker := services.NewBookingWorker()
	worker.Start()

	movieHandler := handlers.NewMovieHandler(movieRepo)
	sessionHandler := handlers.NewSessionHandler(sessionRepo)
	ticketHandler := handlers.NewTicketHandler(ticketRepo, worker)

	routes.RegisterRoutes(movieHandler, ticketHandler, sessionHandler)

	log.Println("Cinema backend running at :8080")
	http.ListenAndServe(":8080", nil)
}
