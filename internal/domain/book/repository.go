package book

// Repository defines the interface for book persistence
type Repository interface {
	Save(book *Book) error
	FindByID(id string) (*Book, error)
	FindByUserID(userID string) ([]*Book, error)
	FindByUserIDAndStatus(userID string, status Status) ([]*Book, error)
	Update(book *Book) error
	Delete(id string) error
}
