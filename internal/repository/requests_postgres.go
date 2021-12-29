package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Dyleme/image-coverter/internal/model"
)

// ReqPostgres is a struct that provide methods get, add, delete and update requests.
type ReqPostgres struct {
	db *sql.DB
}

// NewReqPostgres is constructor for the ReqPostgres.
func NewReqPostgres(db *sql.DB) *ReqPostgres {
	return &ReqPostgres{db: db}
}

// GetRequests method gets all user's requests from the postgres database.
func (r *ReqPostgres) GetRequests(ctx context.Context, userID int) ([]model.Request, error) {
	query := fmt.Sprintf(`SELECT id, op_status, request_time, completion_time, original_id,
	 processed_id, ratio, original_type, processed_type FROM %s WHERE user_id = $1`, RequestTable)

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("repo: %w", err)
	}
	defer rows.Close()

	var reqs []model.Request

	for rows.Next() {
		req := new(model.Request)

		var complTime sql.NullTime

		var processedID sql.NullInt64

		err := rows.Scan(&req.ID, &req.OpStatus, &req.RequestTime, &complTime,
			&req.OriginalID, &processedID, &req.Ratio,
			&req.OriginalType, &req.ProcessedType)

		if err != nil {
			return nil, fmt.Errorf("repo: %w", err)
		}

		if complTime.Valid {
			req.CompletionTime = complTime.Time
		}

		if processedID.Valid {
			req.ProcessedID = int(processedID.Int64)
		}

		reqs = append(reqs, *req)
	}

	return reqs, nil
}

// GetRequests method gets one request from the database by its id.
// If this request belongs to the another user, this function returns error.
func (r *ReqPostgres) GetRequest(ctx context.Context, userID, reqID int) (*model.Request, error) {
	query := fmt.Sprintf(`SELECT id, op_status, request_time, completion_time, original_id,
	 processed_id, ratio, original_type, processed_type FROM %s WHERE id = $1 and user_id = $2`, RequestTable)
	row := r.db.QueryRowContext(ctx, query, reqID, userID)

	var (
		req         model.Request
		complTime   sql.NullTime
		processedID sql.NullInt64
	)

	err := row.Scan(&req.ID, &req.OpStatus, &req.RequestTime, &complTime,
		&req.OriginalID, &processedID, &req.Ratio,
		&req.OriginalType, &req.ProcessedType)
	if err != nil {
		return nil, fmt.Errorf("repo: %w", err)
	}

	if complTime.Valid {
		req.CompletionTime = complTime.Time
	}

	if processedID.Valid {
		req.ProcessedID = int(processedID.Int64)
	}

	return &req, nil
}

// AddRequest method add a request to the database and returns request id.
func (r *ReqPostgres) AddRequest(ctx context.Context, req *model.Request, userID int) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (op_status, request_time, original_id, 
		user_id, ratio, original_type, processed_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`, RequestTable)
	row := r.db.QueryRowContext(ctx, query, req.OpStatus, req.RequestTime, req.OriginalID,
		userID, req.Ratio, req.OriginalType, req.ProcessedType)

	var reqID int
	if err := row.Scan(&reqID); err != nil {
		return 0, fmt.Errorf("repo: %w", err)
	}

	return reqID, nil
}

// AddImage method add image to the postgres database.
// Returns id of this image.
func (r *ReqPostgres) AddImage(ctx context.Context, userID int, imageInfo model.ReuquestImageInfo) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (im_type, image_url, user_id)
		VALUES ($1, $2, $3) RETURNING id;`, ImageTable)
	row := r.db.QueryRowContext(ctx, query, imageInfo.Type, imageInfo.URL, userID)

	var imageID int

	if err := row.Scan(&imageID); err != nil {
		return 0, fmt.Errorf("repo: %w", err)
	}

	return imageID, nil
}

// DeleteRequest method deletes request with reqeust id from database.
// Returns id of the origianal and converted images.
func (r *ReqPostgres) DeleteRequest(ctx context.Context, userID, reqID int) (im1id, im2id int, err error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = $1 AND id = $2 RETURNING original_id, processed_id`, RequestTable)
	row := r.db.QueryRowContext(ctx, query, userID, reqID)

	if err := row.Scan(&im1id, &im2id); err != nil {
		return 0, 0, fmt.Errorf("repo: %w", err)
	}

	return im1id, im2id, nil
}

// DeleteImage method delete image from the database. Returns url path to this image.
func (r *ReqPostgres) DeleteImage(ctx context.Context, userID, imageID int) (string, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = $1 AND id = $2 RETURNING image_url`, ImageTable)
	row := r.db.QueryRowContext(ctx, query, userID, imageID)

	var url string

	if err := row.Scan(&url); err != nil {
		return "", fmt.Errorf("repo: %w", err)
	}

	return url, nil
}
