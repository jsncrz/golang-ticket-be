package router

import (
	"ticket-tracker/api/resource/dashboard"
	"ticket-tracker/api/resource/ticket"
	mwCustom "ticket-tracker/api/router/middleware"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.mongodb.org/mongo-driver/mongo"
)

func New(c *mongo.Client) *chi.Mux {
	r := chi.NewRouter()

	r.Use(mwCustom.CTJson)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://localhost*", "http://localhost*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Route("/v1", func(r chi.Router) {
		ticketAPI := ticket.New(c)
		r.Route("/tickets", func(r chi.Router) {

			r.Get("/", ticketAPI.Get)
			r.Post("/", ticketAPI.Create)
			r.Put("/{id}", ticketAPI.Update)
		})
		r.Route("/calendar", func(r chi.Router) {

			r.Get("/", ticketAPI.GetCalendar)
		})
		dashboardAPI := dashboard.New(c)
		r.Route("/dashboard", func(r chi.Router) {
			r.Get("/tickets/monthly", dashboardAPI.GetMonthlyTicketsCount)
		})
	})
	return r
}
