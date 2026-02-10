package routes

import (
	"cinema-system/internal/handlers"
	"cinema-system/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authHandler    *handlers.AuthHandler
	movieHandler   *handlers.MovieHandler
	sessionHandler *handlers.SessionHandler
	bookingHandler *handlers.BookingHandler
	reviewHandler  *handlers.ReviewHandler
	hallHandler    *handlers.HallHandler
}

func NewRouter(
	authHandler *handlers.AuthHandler,
	movieHandler *handlers.MovieHandler,
	sessionHandler *handlers.SessionHandler,
	bookingHandler *handlers.BookingHandler,
	reviewHandler *handlers.ReviewHandler,
	hallHandler *handlers.HallHandler,
) *Router {
	return &Router{
		authHandler:    authHandler,
		movieHandler:   movieHandler,
		sessionHandler: sessionHandler,
		bookingHandler: bookingHandler,
		reviewHandler:  reviewHandler,
		hallHandler:    hallHandler,
	}
}

func (r *Router) Setup() *gin.Engine {
	router := gin.Default()

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
		user.POST("/bookings", r.bookingHandler.BookTicket)
		user.DELETE("/bookings/:id", r.bookingHandler.CancelTicket)
		user.GET("/bookings/my", r.bookingHandler.GetMyTickets)

		user.POST("/reviews", r.reviewHandler.CreateReview)
		user.GET("/reviews/movie/:movieId", r.reviewHandler.GetMovieReviews)
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
	}

	return router
}
