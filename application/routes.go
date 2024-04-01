package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/paulsheridan/booking-go/handler"
)

func loadRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/clients", loadClientRoutes)

	return router
}

func loadClientRoutes(router chi.Router) {
	clientHandler := &handler.Client{}

	router.Post("/", clientHandler.Create)
	router.Get("/", clientHandler.List)
	router.Get("/{id}", clientHandler.GetByID)
	router.Put("/{id}", clientHandler.UpdateByID)
	router.Delete("/{id}", clientHandler.DeleteByID)
}
