package postgres

import (
	"database/sql"

	"github.com/guisithos/save-my-read/internal/domain/auth"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *auth.User) error {
	query := `
		INSERT INTO users (id, email, name, password, genres)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query,
		user.ID,
		user.Email,
		user.Name,
		user.Password,
		pq.Array(user.Genres),
	)

	return err
}

func (r *UserRepository) FindByEmail(email string) (*auth.User, error) {
	query := `
		SELECT id, email, name, password, genres
		FROM users
		WHERE email = $1
	`

	user := &auth.User{}
	var genres []string

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		pq.Array(&genres),
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	user.Genres = genres
	return user, nil
}

func (r *UserRepository) FindByID(id string) (*auth.User, error) {
	query := `
		SELECT id, email, name, password, genres
		FROM users
		WHERE id = $1
	`

	user := &auth.User{}
	var genres []string

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Password,
		pq.Array(&genres),
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	user.Genres = genres
	return user, nil
}

func (r *UserRepository) Update(user *auth.User) error {
	query := `
		UPDATE users
		SET email = $2, name = $3, password = $4, genres = $5
		WHERE id = $1
	`

	_, err := r.db.Exec(query,
		user.ID,
		user.Email,
		user.Name,
		user.Password,
		pq.Array(user.Genres),
	)

	return err
}
