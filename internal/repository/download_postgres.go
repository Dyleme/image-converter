package repository

import (
	"context"
	"database/sql"
	"fmt"
)

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

	if err := row.Scan(&urlImage); err != nil {
		return "", fmt.Errorf("repo: %w", err)
	}

	return urlImage, nil
}
