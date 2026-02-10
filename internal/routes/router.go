package routes

import (
	"cinema-system/internal/handlers"
	"cinema-system/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler        *handlers.AuthHandler
	movieHandler       *handlers.MovieHandler
	sessionHandler     *handlers.SessionHandler
	bookingHandler     *handlers.BookingHandler
	reviewHandler      *handlers.ReviewHandler
	hallHandler        *handlers.HallHandler
	paymentCardHandler *handlers.PaymentCardHandler
	paymentHandler     *handlers.PaymentHandler
}

func NewRouter(
	authHandler *handlers.AuthHandler,
	movieHandler *handlers.MovieHandler,
	sessionHandler *handlers.SessionHandler,
	bookingHandler *handlers.BookingHandler,
	reviewHandler *handlers.ReviewHandler,
	hallHandler *handlers.HallHandler,
	paymentCardHandler *handlers.PaymentCardHandler,
	paymentHandler *handlers.PaymentHandler,
) *Router {
	return &Router{
		authHandler:        authHandler,
		movieHandler:       movieHandler,
		sessionHandler:     sessionHandler,
		bookingHandler:     bookingHandler,
		reviewHandler:      reviewHandler,
		hallHandler:        hallHandler,
		paymentCardHandler: paymentCardHandler,
		paymentHandler:     paymentHandler,
	}
}

func (r *Router) Setup() *gin.Engine {
	router := gin.Default()
	router.Static("/ui", "./frontend")

	public := router.Group("/api")
	{
		public.POST("/auth/register", r.authHandler.Register)
		public.POST("/auth/login", r.authHandler.Login)

		public.GET("/movies", r.movieHandler.GetAllMovies)
		public.GET("/movies/:id", r.movieHandler.GetMovie)
		public.GET("/sessions/upcoming", r.sessionHandler.GetUpcomingSessions)
		public.GET("/sessions/movie/:movieId", r.sessionHandler.GetMovieSessions)
		public.GET("/sessions/:id", r.sessionHandler.GetSession)
		public.GET("/sessions/:id/booked-seats", r.bookingHandler.GetSessionBookedSeats)
		public.GET("/halls/:id", r.hallHandler.GetHall)
	}

	user := router.Group("/api")
	user.Use(middleware.AuthRequired())
	{
		user.GET("/auth/me", r.authHandler.GetMe)
		user.PUT("/auth/me", r.authHandler.UpdateProfile)
		user.POST("/bookings", r.bookingHandler.BookTickets)
		user.DELETE("/bookings/:id", r.bookingHandler.CancelTicket)
		user.GET("/bookings/my", r.bookingHandler.GetMyTickets)

		user.POST("/reviews", r.reviewHandler.CreateReview)
		user.GET("/reviews/movie/:movieId", r.reviewHandler.GetMovieReviews)
		user.GET("/reviews/my", r.reviewHandler.GetMyReviews)
		user.PUT("/reviews/:id", r.reviewHandler.UpdateReview)
		user.DELETE("/reviews/:id", r.reviewHandler.DeleteReview)

		user.POST("/payment-cards", r.paymentCardHandler.CreateCard)
		user.GET("/payment-cards", r.paymentCardHandler.GetMyCards)
		user.GET("/payment-cards/:id", r.paymentCardHandler.GetCard)
		user.DELETE("/payment-cards/:id", r.paymentCardHandler.DeleteCard)

		user.POST("/payments/topup", r.paymentHandler.TopUpBalance)
		user.GET("/payments", r.paymentHandler.GetMyPayments)
		user.GET("/payments/:id", r.paymentHandler.GetPayment)
		user.POST("/payments/:id/refund", r.paymentHandler.RefundPayment)
	}

	admin := router.Group("/api/admin")
	admin.Use(middleware.AuthRequired(), middleware.AdminRequired())
	{
		admin.POST("/movies", r.movieHandler.CreateMovie)
		admin.PUT("/movies/:id", r.movieHandler.UpdateMovie)
		admin.DELETE("/movies/:id", r.movieHandler.DeleteMovie)

		admin.POST("/halls", r.hallHandler.CreateHall)
		admin.GET("/halls", r.hallHandler.GetAllHalls)
		admin.PUT("/halls/:id", r.hallHandler.UpdateHall)
		admin.DELETE("/halls/:id", r.hallHandler.DeleteHall)

		admin.POST("/sessions", r.sessionHandler.CreateSession)
		admin.DELETE("/sessions/:id", r.sessionHandler.DeleteSession)

		admin.GET("/bookings", r.bookingHandler.GetAllBookings)
		admin.GET("/bookings/session/:sessionId", r.bookingHandler.GetSessionTickets)

		admin.DELETE("/reviews/:id", r.reviewHandler.DeleteReview)

		admin.GET("/payments", r.paymentHandler.GetAllPayments)
		admin.GET("/payments/user/:userId", r.paymentHandler.GetUserPaymentsByID)
		admin.GET("/payment-cards/user/:userId", r.paymentCardHandler.GetUserCards)
	}

	return router
}
