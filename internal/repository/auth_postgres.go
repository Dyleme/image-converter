package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Dyleme/image-coverter/internal/model"
)

// AuthPostgres is a struct that provide methods to create and get information about user from sql.DB.
type AuthPostgres struct {
	db *sql.DB
}

// NewAuthPostgres is a constructor for AuthPostgres.
func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

var ErrUserNotExist = errors.New("such user not exist")

// CreateUser method creates user in database. And returns the id of the user.
func (r *AuthPostgres) CreateUser(ctx context.Context, user model.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (nickname, email, password_hash) VALUES ($1, $2, $3) RETURNING id", UsersTable)
	row := r.db.QueryRowContext(ctx, query, user.Nickname, user.Email, user.Password)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("repo: %w", err)
	}

	return id, nil
}

// GetPasswordHashAndId method returns user id and password hash from the db.
func (r *AuthPostgres) GetPasswordHashAndID(ctx context.Context, nickname string) (hash []byte, userID int, err error) {
	query := fmt.Sprintf("SELECT password_hash, id FROM %s WHERE nickname = $1", UsersTable)
	row := r.db.QueryRowContext(ctx, query, nickname)

	err = row.Scan(&hash, &userID)

	if errors.Is(err, ErrUserNotExist) {
		return nil, 0, fmt.Errorf("repo: %w", ErrUserNotExist)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("repo: %w", err)
	}

	return hash, userID, nil
}
