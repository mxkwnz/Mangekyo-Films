package main

import (
	"cinema-system/internal/config"
	"cinema-system/internal/handlers"
	"cinema-system/internal/repositories"
	"cinema-system/internal/routes"
	"cinema-system/internal/services"
	"log"
	"os"
)

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	db, err := config.NewDatabase(mongoURI, "cinema_db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Disconnect()

	userRepo := repositories.NewUserRepository(db.Database)
	movieRepo := repositories.NewMovieRepository(db.Database)
	genreRepo := repositories.NewGenreRepository(db.Database)
	hallRepo := repositories.NewHallRepository(db.Database)
	sessionRepo := repositories.NewSessionRepository(db.Database)
	ticketRepo := repositories.NewTicketRepository(db.Database)
	reviewRepo := repositories.NewReviewRepository(db.Database)

	authService := services.NewAuthService(userRepo)
	movieService := services.NewMovieService(movieRepo, genreRepo)
	sessionService := services.NewSessionService(sessionRepo, hallRepo, movieRepo)
	bookingService := services.NewBookingService(ticketRepo, sessionRepo, userRepo, hallRepo)
	reviewService := services.NewReviewService(reviewRepo, movieRepo, userRepo)

	authHandler := handlers.NewAuthHandler(authService)
	movieHandler := handlers.NewMovieHandler(movieService)
	sessionHandler := handlers.NewSessionHandler(sessionService)
	bookingHandler := handlers.NewBookingHandler(bookingService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	hallHandler := handlers.NewHallHandler(hallRepo)

	router := routes.NewRouter(
		authHandler,
		movieHandler,
		sessionHandler,
		bookingHandler,
		reviewHandler,
		hallHandler,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Cinema System Server starting on port %s", port)
	if err := router.Setup().Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
