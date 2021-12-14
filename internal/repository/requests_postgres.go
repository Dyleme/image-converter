package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Dyleme/image-coverter/internal/model"
)

type ReqPostgres struct {
	db *sql.DB
}

func NewReqPostgres(db *sql.DB) *ReqPostgres {
	return &ReqPostgres{db: db}
}

func (r *ReqPostgres) GetRequests(ctx context.Context, userID int) ([]model.Request, error) {
	query := fmt.Sprintf(`SELECT id, op_status, request_time, completion_time, original_id,
	 processed_id, ratio, original_type, processed_type FROM %s WHERE user_id = $1`, RequestTable)

	rows, err := r.db.Query(query, userID)

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

func (r *ReqPostgres) GetRequest(ctx context.Context, userID, reqID int) (*model.Request, error) {
	query := fmt.Sprintf(`SELECT id, op_status, request_time, completion_time, original_id,
	 processed_id, ratio, original_type, processed_type FROM %s WHERE id = $1 and user_id = $2`, RequestTable)

	row := r.db.QueryRow(query, reqID, userID)

	var req model.Request

	var complTime sql.NullTime

	var processedID sql.NullInt64

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

func (r *ReqPostgres) AddRequest(ctx context.Context, req *model.Request, userID int) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (op_status, request_time, original_id, 
		user_id, ratio, original_type, processed_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`, RequestTable)
	row := r.db.QueryRow(query, req.OpStatus, req.RequestTime, req.OriginalID,
		userID, req.Ratio, req.OriginalType, req.ProcessedType)

	var reqID int
	if err := row.Scan(&reqID); err != nil {
		return 0, fmt.Errorf("repo: %w", err)
	}

	return reqID, nil
}

func (r *ReqPostgres) UpdateRequestStatus(ctx context.Context, reqID int, status string) error {
	query := fmt.Sprintf(`UPDATE %s SET op_status = $1 WHERE id = $2 RETURNING id`, RequestTable)
	row := r.db.QueryRow(query, status, reqID)

	var id int
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

func (r *ReqPostgres) AddProcessedImageIDToRequest(ctx context.Context, reqID, imageID int) error {
	query := fmt.Sprintf(`UPDATE %s SET processed_id = $1 WHERE id = $2 RETURNING id;`, RequestTable)
	row := r.db.QueryRow(query, imageID, reqID)

	var id int
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

func (r *ReqPostgres) AddProcessedTimeToRequest(ctx context.Context, reqID int, t time.Time) error {
	query := fmt.Sprintf(`UPDATE %s SET completion_time = $1 WHERE id = $2 RETURNING id;`, RequestTable)
	row := r.db.QueryRow(query, t, reqID)

	var id int
	if err := row.Scan(&id); err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

func (r *ReqPostgres) AddImage(ctx context.Context, userID int, imageInfo model.Info) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (resoolution_x, resoolution_y, im_type, image_url, user_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id;`, ImageTable)
	row := r.db.QueryRow(query, imageInfo.Width, imageInfo.Height,
		imageInfo.Type, imageInfo.URL, userID)

	var imageID int
	if err := row.Scan(&imageID); err != nil {
		return 0, fmt.Errorf("repo: %w", err)
	}

	return imageID, nil
}

func (r *ReqPostgres) DeleteRequest(ctx context.Context, userID, reqID int) (im1id, im2id int, err error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = $1 AND id = $2 RETURNING original_id, processed_id`, RequestTable)

	row := r.db.QueryRow(query, userID, reqID)

	err = row.Scan(&im1id, &im2id)

	if err != nil {
		return 0, 0, fmt.Errorf("repo: %w", err)
	}

	return im1id, im2id, nil
}

func (r *ReqPostgres) DeleteImage(ctx context.Context, userID, imageID int) (string, error) {
	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = $1 AND id = $2 RETURNING image_url`, ImageTable)

	row := r.db.QueryRow(query, userID, imageID)

	var url string

	if err := row.Scan(&url); err != nil {
		return "", fmt.Errorf("repo: %w", err)
	}

	return url, nil
}
