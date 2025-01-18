package server

import (
	"log"
	"net/http"

	"github.com/guisithos/save-my-read/internal/interfaces/http/handlers"
	"github.com/guisithos/save-my-read/internal/interfaces/http/middleware"
)

// Server represents the HTTP server
type Server struct {
	bookHandler *handlers.BookHandler
	authHandler *handlers.AuthHandler
	viewHandler *handlers.ViewHandler
	port        string
	jwtKey      []byte
}

// NewServer creates a new HTTP server
func NewServer(bookHandler *handlers.BookHandler, authHandler *handlers.AuthHandler, port string, jwtKey string) *Server {
	viewHandler, err := handlers.NewViewHandler()
	if err != nil {
		log.Fatalf("Failed to create view handler: %v", err)
	}

	return &Server{
		bookHandler: bookHandler,
		authHandler: authHandler,
		viewHandler: viewHandler,
		port:        port,
		jwtKey:      []byte(jwtKey),
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Static files
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Create a base handler with CORS middleware
	baseHandler := addMiddleware(mux)

	// View routes
	mux.HandleFunc("/", s.viewHandler.Home)

	// Auth routes (no auth middleware)
	mux.HandleFunc("/api/auth/register", addMiddleware(http.HandlerFunc(s.authHandler.Register)).ServeHTTP)
	mux.HandleFunc("/api/auth/login", addMiddleware(http.HandlerFunc(s.authHandler.Login)).ServeHTTP)

	// Protected API routes (with auth middleware)
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/books/search", s.bookHandler.SearchBooks)
	protectedMux.HandleFunc("/api/books/add", s.bookHandler.AddToList)
	protectedMux.HandleFunc("/api/books/user", s.bookHandler.GetUserBooks)
	protectedMux.HandleFunc("/api/books/status", s.bookHandler.UpdateBookStatus)

	// Add auth middleware only to protected routes
	protectedHandler := middleware.AuthMiddleware(s.jwtKey)(addMiddleware(protectedMux))

	// Combine handlers
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || r.URL.Path == "/api/auth/register" || r.URL.Path == "/api/auth/login" {
			baseHandler.ServeHTTP(w, r)
			return
		}
		protectedHandler.ServeHTTP(w, r)
	})

	log.Printf("Server starting on port %s", s.port)
	return http.ListenAndServe(":"+s.port, finalHandler)
}

func addMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Log request
		log.Printf("%s %s", r.Method, r.URL.Path)

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
