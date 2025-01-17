package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/guisithos/save-my-read/internal/application"
	"github.com/guisithos/save-my-read/internal/domain/auth"
	"github.com/guisithos/save-my-read/internal/infrastructure/postgres"
	"github.com/guisithos/save-my-read/internal/interfaces/http/handlers"
	"github.com/guisithos/save-my-read/internal/interfaces/http/server"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/save_my_read?sslmode=disable"
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)

	// Initialize services
	jwtService := auth.NewJWTService(os.Getenv("JWT_SECRET"), 24*time.Hour)
	authService := application.NewAuthService(userRepo, jwtService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize and start server
	srv := server.NewServer(authHandler)
	log.Fatal(http.ListenAndServe(":8080", srv.Router()))
}
