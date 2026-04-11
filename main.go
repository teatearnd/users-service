package main

import (
	"log"
	"net/http"

	"example.com/internal/handlers"
	"example.com/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	if err := auth.ValidateConfig(); err != nil {
		log.Fatalf("JWT_SECRET not configured properly in the .env: %w", err)
	}

	r.Post("/login", handlers.LoginHandler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Println("starting at :8081")
	if err := http.ListenAndServe(":8081", r); err != nil {
		log.Fatalf("server is down: %v", err)
	}
}
