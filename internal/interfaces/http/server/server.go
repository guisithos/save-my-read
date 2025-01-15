package server

import (
	"log"
	"net/http"

	"github.com/guisithos/save-my-read/internal/interfaces/http/handlers"
)

// Server represents the HTTP server
type Server struct {
	bookHandler *handlers.BookHandler
	viewHandler *handlers.ViewHandler
	port        string
}

// NewServer creates a new HTTP server
func NewServer(bookHandler *handlers.BookHandler, port string) *Server {
	viewHandler, err := handlers.NewViewHandler()
	if err != nil {
		log.Fatalf("Failed to create view handler: %v", err)
	}

	return &Server{
		bookHandler: bookHandler,
		viewHandler: viewHandler,
		port:        port,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// Static files
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// View routes
	mux.HandleFunc("/", s.viewHandler.Home)

	// API routes
	mux.HandleFunc("/api/books/search", s.bookHandler.SearchBooks)
	mux.HandleFunc("/api/books/add", s.bookHandler.AddToList)
	mux.HandleFunc("/api/books/user", s.bookHandler.GetUserBooks)
	mux.HandleFunc("/api/books/status", s.bookHandler.UpdateBookStatus)

	// Add middleware
	handler := addMiddleware(mux)

	log.Printf("Server starting on port %s", s.port)
	return http.ListenAndServe(":"+s.port, handler)
}

func addMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

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
