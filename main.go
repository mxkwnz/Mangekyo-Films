package main

import (
	"cinema-system/internal/config"
	"cinema-system/internal/handlers"
	"cinema-system/internal/repositories"
	"cinema-system/internal/routes"
	"cinema-system/internal/services"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGO_DATABASE")
	if dbName == "" {
		dbName = "cinema_db"
	}

	db, err := config.NewDatabase(mongoURI, dbName)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Disconnect()

	log.Println("Successfully connected to MongoDB!")

	userRepo := repositories.NewUserRepository(db.Database)
	movieRepo := repositories.NewMovieRepository(db.Database)
	genreRepo := repositories.NewGenreRepository(db.Database)
	hallRepo := repositories.NewHallRepository(db.Database)
	sessionRepo := repositories.NewSessionRepository(db.Database)
	ticketRepo := repositories.NewTicketRepository(db.Database)
	reviewRepo := repositories.NewReviewRepository(db.Database)
	paymentCardRepo := repositories.NewPaymentCardRepository(db.Database)
	paymentRepo := repositories.NewPaymentRepository(db.Database)

	authService := services.NewAuthService(userRepo)
	movieService := services.NewMovieService(movieRepo, genreRepo)
	sessionService := services.NewSessionService(sessionRepo, hallRepo, movieRepo)
	bookingService := services.NewBookingService(ticketRepo, sessionRepo, userRepo, hallRepo, paymentRepo)
	reviewService := services.NewReviewService(reviewRepo, movieRepo, userRepo)
	paymentCardService := services.NewPaymentCardService(paymentCardRepo, userRepo)
	paymentService := services.NewPaymentService(paymentRepo, paymentCardRepo, userRepo)

	authHandler := handlers.NewAuthHandler(authService)
	movieHandler := handlers.NewMovieHandler(movieService)
	sessionHandler := handlers.NewSessionHandler(sessionService)
	bookingHandler := handlers.NewBookingHandler(bookingService)
	reviewHandler := handlers.NewReviewHandler(reviewService)
	hallHandler := handlers.NewHallHandler(hallRepo)
	paymentCardHandler := handlers.NewPaymentCardHandler(paymentCardService)
	paymentHandler := handlers.NewPaymentHandler(paymentService)

	router := routes.NewRouter(
		authHandler,
		movieHandler,
		sessionHandler,
		bookingHandler,
		reviewHandler,
		hallHandler,
		paymentCardHandler,
		paymentHandler,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🎬 Cinema System Server starting on port %s", port)
	if err := router.Setup().Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
