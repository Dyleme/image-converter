package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Dyleme/image-coverter/internal/model"
)

// ConvPostgres is a struct that provides methods to add image and update it's resolution int the sql.DB.
type ConvPostgres struct {
	db *TxDB
}

// NewConvPostgres is a constructor for the ConvPostgers.
func NewConvPostgres(db *sql.DB) *ConvPostgres {
	return &ConvPostgres{db: &TxDB{db}}
}

// GetConvInfo method returns all information about request from database.
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
		return nil, err
	}

	return &inf, nil
}

// SetImageResolution method set image resolution to the image in images table.
func (c *ConvPostgres) SetImageResolution(ctx context.Context, imID, width, height int) error {
	query := fmt.Sprintf(`UPDATE %s 
	SET resoolution_x = $1,
	resoolution_y = $2
	WHERE id = $3 
	`, ImageTable)

	result, err := c.db.ExecContext(ctx, query, width, height, imID)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return oneRowInResult(result)
}

// AddImageDB is a colmplex method that creates thransaction.
// And in this transaction at first it add image to the images table.
// Then it sets resolution of this image. After it add this image, processed time
// to the requests table and updates request status.
// Returns any error occurred in transaction or while creatring transaction.
func (c *ConvPostgres) AddImageDB(ctx context.Context, userID, reqID int, imgInfo *model.ReuquestImageInfo,
	width, height int, status string, t time.Time) error {
	err := c.db.inTx(ctx, func(tx *sql.Tx) error {
		imID, err := addImageWithResolution(ctx, tx, userID, *imgInfo, width, height)
		if err != nil {
			return err
		}

		err = addProcessedImageIDToRequest(ctx, tx, reqID, imID)
		if err != nil {
			return err
		}

		err = addProcessedTimeToRequest(ctx, tx, reqID, t)
		if err != nil {
			return err
		}

		err = updateRequestStatus(ctx, tx, reqID, status)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

// updateRequestStatus function update status of an existing request in database.
func updateRequestStatus(ctx context.Context, tx *sql.Tx, reqID int, status string) error {
	query := fmt.Sprintf(`UPDATE %s SET op_status = $1 WHERE id = $2;`, RequestTable)

	result, err := tx.ExecContext(ctx, query, status, reqID)
	if err != nil {
		return err
	}

	return oneRowInResult(result)
}

// addProcessedImageIDToRequest function update processed image id column for the reqId.
func addProcessedImageIDToRequest(ctx context.Context, tx *sql.Tx, reqID, imageID int) error {
	query := fmt.Sprintf(`UPDATE %s SET processed_id = $1 WHERE id = $2;`, RequestTable)

	result, err := tx.ExecContext(ctx, query, imageID, reqID)
	if err != nil {
		return err
	}

	return oneRowInResult(result)
}

// addImageToDB function add image to the postgres database.
// Returns id of this image.
func addImageWithResolution(ctx context.Context, tx *sql.Tx, userID int,
	imageInfo model.ReuquestImageInfo, width, height int) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (im_type, image_url, user_id, resoolution_x, resoolution_y)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`, ImageTable)
	row := tx.QueryRowContext(ctx, query, imageInfo.Type, imageInfo.URL, userID, width, height)

	var imageID int
	if err := row.Scan(&imageID); err != nil {
		return 0, err
	}

	return imageID, nil
}

// AddProcessedImageIDToRequest method update processed time column for the reqId.
func addProcessedTimeToRequest(ctx context.Context, tx *sql.Tx, reqID int, t time.Time) error {
	query := fmt.Sprintf(`UPDATE %s SET completion_time = $1 WHERE id = $2;`, RequestTable)

	result, err := tx.ExecContext(ctx, query, t, reqID)
	if err != nil {
		return err
	}

	return oneRowInResult(result)
}
