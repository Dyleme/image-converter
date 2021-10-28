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

func (r *ReqPostgres) GetRequests(id int) ([]image.Request, error) {
	query := fmt.Sprintf(`SELECT id, op_status, request_time, completion_time, original_id,
	 processed_id, ratio, original_type, processed_type FROM %s WHERE user_id = $1`, requestTable)

	rows, err := r.db.Query(query, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reqs []image.Request

	fmt.Println(id)

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
