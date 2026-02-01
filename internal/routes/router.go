package routes

import (
	"cinema-system/internal/handlers"
	"net/http"
)

func RegisterRoutes(
	movieHandler *handlers.MovieHandler,
) {
	http.HandleFunc("/movies", movieHandler.HandleMovies)
}
