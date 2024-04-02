package application

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/paulsheridan/booking-go/database/client"
	"github.com/paulsheridan/booking-go/handler"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/clients", a.loadClientRoutes)

	a.router = router
}

func (a *App) loadClientRoutes(router chi.Router) {
	clientHandler := &handler.Client{
		Repo: &client.RedisDatabase{
			Client: a.rdb,
		},
	}

	router.Post("/", clientHandler.Create)
	router.Get("/", clientHandler.List)
	router.Get("/{id}", clientHandler.GetByID)
	router.Put("/{id}", clientHandler.UpdateByID)
	router.Delete("/{id}", clientHandler.DeleteByID)
}
