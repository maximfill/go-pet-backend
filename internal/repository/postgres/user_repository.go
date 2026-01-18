package postgres

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(
	ctx context.Context,
	email string,
	passwordHash string,
) (int64, error) {
	var id int64

	err := r.db.QueryRowContext(
		ctx,
		`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id
		`,
		email,
		passwordHash,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := User{}

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT id, email, password_hash, created_at
		FROM users
		WHERE email = $1
		`,
		email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.PasswordHash,
		&u.CreatedAt,
	)

	if err != nil {
		return nil, err
	}
	return &u, nil
}
