package server

import (
	"log"
	"net/http"

	"github.com/guisithos/save-my-read/internal/domain/auth"
	"github.com/guisithos/save-my-read/internal/interfaces/http/handlers"
	"github.com/guisithos/save-my-read/internal/interfaces/http/middleware"
)

// Server represents the HTTP server
type Server struct {
	authHandler  *handlers.AuthHandler
	bookHandler  *handlers.BookHandler
	tokenService auth.TokenService
	port         string
}

// NewServer creates a new HTTP server
func NewServer(authHandler *handlers.AuthHandler, bookHandler *handlers.BookHandler, tokenService auth.TokenService, port string) *Server {
	return &Server{
		authHandler:  authHandler,
		bookHandler:  bookHandler,
		tokenService: tokenService,
		port:         port,
	}
}

// SetupRoutes sets up the routes for the HTTP server
func (s *Server) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Public routes (no auth required)
	mux.HandleFunc("/api/auth/register", s.authHandler.Register)
	mux.HandleFunc("/api/auth/login", s.authHandler.Login)
	mux.HandleFunc("/api/books/search", s.bookHandler.SearchBooks)

	// Protected routes (auth required)
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/books", s.bookHandler.GetBooks)
	protectedMux.HandleFunc("/api/books/add", s.bookHandler.AddBook)
	protectedMux.HandleFunc("/api/books/status", s.bookHandler.UpdateBookStatus)

	// Apply auth middleware to protected routes
	authMiddleware := middleware.NewAuthMiddleware(s.tokenService)
	mux.Handle("/api/books/", authMiddleware(protectedMux))

	// Serve static files and templates
	fs := http.FileServer(http.Dir("web/templates"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "web/templates/index.html")
			return
		}
		fs.ServeHTTP(w, r)
	})

	return mux
}

// Start starts the HTTP server
func (s *Server) Start() error {
	mux := s.SetupRoutes()

	log.Printf("Server starting on port %s", s.port)
	return http.ListenAndServe(":"+s.port, mux)
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
