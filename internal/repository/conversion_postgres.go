package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Dyleme/image-coverter/internal/model"
)

type ConvPostgres struct {
	db *sql.DB
}

func NewConvPostgres(db *sql.DB) *ConvPostgres {
	return &ConvPostgres{db: db}
}

func (c *ConvPostgres) GetConvInfo(ctx context.Context, reqID int) (*model.ConvImageInfo, error) {
	query := fmt.Sprintf(`SELECT 
r.user_id, r.original_id, i.image_url, r.original_type, r.processed_type, r.ratio
FROM
%s as r
INNER JOIN 
%s as i
	ON r.original_id = i.id
WHERE 
	r.id = $1`, RequestTable, ImageTable)

	row := c.db.QueryRowContext(ctx, query, reqID)

	var inf model.ConvImageInfo

	err := row.Scan(&inf.UserID, &inf.OldImID, &inf.OldURL, &inf.OldType, &inf.NewType, &inf.Ratio)
	if err != nil {
		return nil, fmt.Errorf("repo: %w", err)
	}

	return &inf, nil
}

// UpdateRequestStatus method update status of an existing request in database.
func (c *ConvPostgres) UpdateRequestStatus(ctx context.Context, reqID int, status string) error {
	query := fmt.Sprintf(`UPDATE %s SET op_status = $1 WHERE id = $2 RETURNING id;`, RequestTable)
	row := c.db.QueryRowContext(ctx, query, status, reqID)

	var id int
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

func (c *ConvPostgres) SetImageResolution(ctx context.Context, imID, widith, height int) error {
	query := fmt.Sprintf(`UPDATE %s 
	SET resoolution_x = $1,
	resoolution_y = $2
	WHERE id = $3 
	RETURNING id;`, ImageTable)
	row := c.db.QueryRowContext(ctx, query, widith, height, imID)

	var id int
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

// AddProcessedImageIDToRequest method update processed image id column for the reqId.
func (c *ConvPostgres) AddProcessedImageIDToRequest(ctx context.Context, reqID, imageID int) error {
	query := fmt.Sprintf(`UPDATE %s SET processed_id = $1 WHERE id = $2 RETURNING id;`, RequestTable)
	row := c.db.QueryRowContext(ctx, query, imageID, reqID)

	var id int
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

// AddProcessedImageIDToRequest method update processed time column for the reqId.
func (c *ConvPostgres) AddProcessedTimeToRequest(ctx context.Context, reqID int, t time.Time) error {
	query := fmt.Sprintf(`UPDATE %s SET completion_time = $1 WHERE id = $2 RETURNING id;`, RequestTable)
	row := c.db.QueryRowContext(ctx, query, t, reqID)

	var id int
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

// AddImage method add image to the postgres database.
// Returns id of this image.
func (c *ConvPostgres) AddImage(ctx context.Context, userID int, imageInfo model.ReuquestImageInfo) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (im_type, image_url, user_id)
		VALUES ($1, $2, $3) RETURNING id;`, ImageTable)
	row := c.db.QueryRowContext(ctx, query, imageInfo.Type, imageInfo.URL, userID)

	var imageID int

	if err := row.Scan(&imageID); err != nil {
		return 0, fmt.Errorf("repo: %w", err)
	}

	return imageID, nil
}
