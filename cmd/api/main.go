package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/guisithos/save-my-read/internal/application"
	"github.com/guisithos/save-my-read/internal/domain/auth"
	"github.com/guisithos/save-my-read/internal/infrastructure/googlebooks"
	"github.com/guisithos/save-my-read/internal/infrastructure/postgres"
	"github.com/guisithos/save-my-read/internal/interfaces/http/handlers"
	"github.com/guisithos/save-my-read/internal/interfaces/http/server"
)

func main() {
	// Verify environment variables
	if os.Getenv("GOOGLE_BOOKS_API_KEY") == "" {
		log.Fatal("GOOGLE_BOOKS_API_KEY environment variable is not set")
	}
	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	// Initialize database connection
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Initialize repositories
	bookRepo := postgres.NewBookRepository(db)
	userRepo := postgres.NewUserRepository(db)

	// Initialize Google Books client
	googleClient, err := googlebooks.NewClient()
	if err != nil {
		log.Fatal("Failed to create Google Books client:", err)
	}

	// Initialize JWT service with 24h token duration
	jwtService := auth.NewJWTService(os.Getenv("JWT_SECRET"), 24*time.Hour)

	// Initialize services
	bookService := application.NewBookService(bookRepo, userRepo)
	authService := application.NewAuthService(userRepo, jwtService)

	// Initialize handlers
	bookHandler := handlers.NewBookHandler(bookService, googleClient)
	authHandler := handlers.NewAuthHandler(authService)

	// Initialize and start server
	srv := server.NewServer(authHandler, bookHandler, jwtService, "8080")
	log.Fatal(srv.Start())
}
