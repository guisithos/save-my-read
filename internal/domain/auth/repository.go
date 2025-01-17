package auth

// UserRepository defines the interface for user persistence operations
type UserRepository interface {
	Create(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
	Update(user *User) error
}
