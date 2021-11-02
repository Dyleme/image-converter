package repository

import (
	"database/sql"
	"fmt"

	"github.com/Dyleme/image-coverter"
)

type ReqPostgres struct {
	db *sql.DB
}

func NewReqPostgres(db *sql.DB) *ReqPostgres {
	return &ReqPostgres{db: db}
}

func (r *ReqPostgres) GetRequests(userID int) ([]image.Request, error) {
	query := fmt.Sprintf(`SELECT id, op_status, request_time, completion_time, original_id,
	 processed_id, ratio, original_type, processed_type FROM %s WHERE user_id = $1`, requestTable)

	rows, err := r.db.Query(query, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reqs []image.Request

	fmt.Println(userID)

	for rows.Next() {
		req := new(image.Request)

		var complTime sql.NullTime

		var processedID sql.NullInt64

		err := rows.Scan(&req.ID, &req.OpStatus, &req.RequestTime, &complTime,
			&req.OriginalID, &processedID, &req.Ratio,
			&req.OriginalType, &req.ProcessedType)

		if err != nil {
			return nil, err
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

func (r *ReqPostgres) GetRequest(userID, reqID int) (*image.Request, error) {
	query := fmt.Sprintf(`SELECT id, op_status, request_time, completion_time, original_id,
	 processed_id, ratio, original_type, processed_type FROM %s WHERE id = $1 and user_id = $2`, requestTable)

	row := r.db.QueryRow(query, reqID, userID)

	var req image.Request

	var complTime sql.NullTime

	var processedID sql.NullInt64

	err := row.Scan(&req.ID, &req.OpStatus, &req.RequestTime, &complTime,
		&req.OriginalID, &processedID, &req.Ratio,
		&req.OriginalType, &req.ProcessedType)

	if err != nil {
		return nil, err
	}

	if complTime.Valid {
		req.CompletionTime = complTime.Time
	}

	if processedID.Valid {
		req.ProcessedID = int(processedID.Int64)
	}

	return &req, nil
}

func (r *ReqPostgres) AddRequest(req *image.Request, userID int) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (op_status, request_time, original_id, 
		user_id, ratio, original_type, processed_type)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`, requestTable)
	row := r.db.QueryRow(query, req.OpStatus, req.RequestTime, req.OriginalID,
		userID, req.Ratio, req.OriginalType, req.ProcessedType)

	var reqID int
	if err := row.Scan(&reqID); err != nil {
		return 0, err
	}

	return reqID, nil
}

func (r *ReqPostgres) AddImage(userID int, imageInfo image.Info) (int, error) {
	query := fmt.Sprintf(`INSERT INTO %s (resoolution_x, resoolution_y, im_type, image_url, user_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id;`, imageTable)
	row := r.db.QueryRow(query, imageInfo.ResoultionX, imageInfo.ResoultionY,
		imageInfo.Type, imageInfo.URL, userID)

	var imageID int
	if err := row.Scan(&imageID); err != nil {
		return 0, err
	}

	return imageID, nil
}
