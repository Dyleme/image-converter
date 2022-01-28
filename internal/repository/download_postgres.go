package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type ImageNotExistError struct{}

func (ImageNotExistError) Error() string {
	return "such image not exist"
}

func (ImageNotExistError) Subject() string {
	return "image"
}

// DownloadPostgres is a struct that provide method to get image url from the sql.DB.
type DownloadPostgres struct {
	db *sql.DB
}

// NewDonwloadPostgres is a constructor for DownloadPostgres.
func NewDownloadPostgres(db *sql.DB) *DownloadPostgres {
	return &DownloadPostgres{db: db}
}

// GetImageUrl function get image url from the database.
func (d *DownloadPostgres) GetImageURL(ctx context.Context, userID, imageID int) (string, error) {
	query := fmt.Sprintf(`SELECT image_url FROM %s WHERE user_id = $1 AND id = $2`, ImageTable)
	row := d.db.QueryRowContext(ctx, query, userID, imageID)

	var urlImage string

	err := row.Scan(&urlImage)

	if errors.Is(err, sql.ErrNoRows) {
		return "", fmt.Errorf("repo: %w", ImageNotExistError{})
	}

	if err != nil {
		return "", fmt.Errorf("repo: %w", err)
	}

	return urlImage, nil
}
