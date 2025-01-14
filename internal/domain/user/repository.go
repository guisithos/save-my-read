package user

// Repository defines the interface for user persistence
type Repository interface {
	Save(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id string) (*User, error)
	Update(user *User) error
}
