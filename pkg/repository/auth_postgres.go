package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Dyleme/image-coverter/pkg/model"
)

type AuthPostgres struct {
	db *sql.DB
}

func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(ctx context.Context, user model.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (nickname, email, password_hash) VALUES ($1, $2, $3) RETURNING id", UsersTable)
	row := r.db.QueryRow(query, user.Nickname, user.Email, user.Password)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("repo: %w", err)
	}

	return id, nil
}

func (r *AuthPostgres) GetPasswordAndID(ctx context.Context, nickname string) (hash []byte, userID int, err error) {
	query := fmt.Sprintf("SELECT password_hash, id FROM %s WHERE nickname = $1", UsersTable)
	row := r.db.QueryRow(query, nickname)

	if row == nil {
		return nil, 0, fmt.Errorf("repo: %w", err)
	}

	if err := row.Scan(&hash, &userID); err != nil {
		return nil, 0, fmt.Errorf("repo: %w", err)
	}

	return hash, userID, nil
}
