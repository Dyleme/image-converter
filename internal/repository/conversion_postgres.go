package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Dyleme/image-coverter/internal/model"
)

type NotSingleRowAffectedError struct {
	amountAffected int
}

func (e *NotSingleRowAffectedError) Error() string {
	return fmt.Sprintf("expected single row affected, got %v rows affected", e.amountAffected)
}

type ConvPostgres struct {
	db *TxDB
}

func NewConvPostgres(db *sql.DB) *ConvPostgres {
	return &ConvPostgres{db: &TxDB{db}}
}

func (c *ConvPostgres) GetConvInfo(ctx context.Context, reqID int) (*model.ConvImageInfo, error) {
	var inf model.ConvImageInfo

	err := c.db.inTx(ctx, func(tx *sql.Tx) error {
		query := fmt.Sprintf(`SELECT 
r.user_id, r.original_id, i.image_url, r.original_type, r.processed_type, r.ratio
FROM
%s as r
INNER JOIN 
%s as i
	ON r.original_id = i.id
WHERE 
	r.id = $1`, RequestTable, ImageTable)

		row := tx.QueryRowContext(ctx, query, reqID)

		return row.Scan(&inf.UserID, &inf.OldImID, &inf.OldURL, &inf.OldType, &inf.NewType, &inf.Ratio)
	})
	if err != nil {
		return nil, fmt.Errorf("repo: %w", err)
	}

	return &inf, nil
}

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

func setImageResolution(ctx context.Context, tx *sql.Tx, imID, width, height int) error {
	query := fmt.Sprintf(`UPDATE %s 
	SET resoolution_x = $1,
	resoolution_y = $2
	WHERE id = $3 
	`, ImageTable)

	result, err := tx.ExecContext(ctx, query, width, height, imID)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return oneRowInResult(result)
}

func (c *ConvPostgres) AddImageDB(ctx context.Context, userID, reqID int, imgInfo *model.ReuquestImageInfo,
	width, height int, status string, t time.Time) error {
	err := c.db.inTx(ctx, func(tx *sql.Tx) error {
		imID, err := addImageToDB(ctx, tx, userID, *imgInfo)
		if err != nil {
			return err
		}

		err = setImageResolution(ctx, tx, imID, width, height)
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

	return err
}

// UpdateRequestStatus method update status of an existing request in database.
func updateRequestStatus(ctx context.Context, tx *sql.Tx, reqID int, status string) error {
	query := fmt.Sprintf(`UPDATE %s SET op_status = $1 WHERE id = $2;`, RequestTable)

	result, err := tx.ExecContext(ctx, query, status, reqID)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return oneRowInResult(result)
}

// AddProcessedImageIDToRequest method update processed image id column for the reqId.
func addProcessedImageIDToRequest(ctx context.Context, tx *sql.Tx, reqID, imageID int) error {
	query := fmt.Sprintf(`UPDATE %s SET processed_id = $1 WHERE id = $2;`, RequestTable)

	result, err := tx.ExecContext(ctx, query, imageID, reqID)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return oneRowInResult(result)
}

// AddImage method add image to the postgres database.
// Returns id of this image.
func addImageToDB(ctx context.Context, tx *sql.Tx, userID int, imageInfo model.ReuquestImageInfo) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (im_type, image_url, user_id)
		VALUES ($1, $2, $3) RETURNING id;`, ImageTable)
	row := tx.QueryRowContext(ctx, query, imageInfo.Type, imageInfo.URL, userID)

	var imageID int
	if err := row.Scan(&imageID); err != nil {
		return 0, fmt.Errorf("repo: %w", err)
	}

	return imageID, nil
}

// AddProcessedImageIDToRequest method update processed time column for the reqId.
func addProcessedTimeToRequest(ctx context.Context, tx *sql.Tx, reqID int, t time.Time) error {
	query := fmt.Sprintf(`UPDATE %s SET completion_time = $1 WHERE id = $2;`, RequestTable)

	result, err := tx.ExecContext(ctx, query, t, reqID)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return oneRowInResult(result)
}
