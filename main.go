package main

import (
	"log"
	"net/http"

	"example.com/internal/config"
	"example.com/internal/handlers"
	"example.com/internal/repository"
	"example.com/pkg/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf(".env not loaded: %v", err)
	}

	cfg := config.LoadConfig()
	authInit := auth.Settings{
		Secret:   cfg.JWTSecret,
		Issuer:   cfg.JWTIssuer,
		Audience: cfg.JWTAudience,
	}

	if err := auth.Init(authInit); err != nil {
		log.Fatalf("JWT secret init failed: %v", err)
	}

	db, err := repository.OpenDB(*cfg)
	if err != nil {
		log.Fatalf("failed at db open: %v", err)
	}
	defer db.Close()
	err = repository.InitSchema(db)
	if err != nil {
		log.Fatalf("failed at db init: %v", err)
	}

	def_handler := &handlers.Handler{DB: db}

	r := chi.NewRouter()
	r.Use(middleware.Logger) // todo make an auth middleware for future routes (except login/register)
	if err := auth.ValidateConfig(); err != nil {
		log.Fatalf("JWT_SECRET not configured properly in the .env: %v", err)
	}

	r.Post("/login", def_handler.LoginHandler)
	r.Post("/register", def_handler.RegisterHandler)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("starting at %s", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, r); err != nil {
		log.Fatalf("server is down: %v", err) // expects a ":port", not a port
	}
}
