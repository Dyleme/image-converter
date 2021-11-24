package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Dyleme/image-coverter/pkg/model"
	"github.com/Dyleme/image-coverter/pkg/repository"
	"github.com/google/go-cmp/cmp"
)

func TestGetRequest(t *testing.T) {
	testCases := []struct {
		testName string
		userID   int
		reqID    int
		repoReq  *model.Request
		wantReq  *model.Request
		wantErr  error
	}{
		{
			testName: "all is good",
			userID:   12,
			reqID:    19,
			repoReq: &model.Request{
				ID:             24,
				OpStatus:       "done",
				RequestTime:    time.Date(2020, 12, 12, 23, 23, 0, 1, time.Local),
				CompletionTime: time.Date(2020, 12, 12, 23, 24, 0, 1, time.Local),
				OriginalID:     12,
				ProcessedID:    13,
				Ratio:          0.5,
				OriginalType:   "jpeg",
				ProcessedType:  "png",
			},
			wantReq: &model.Request{
				ID:             24,
				OpStatus:       "done",
				RequestTime:    time.Date(2020, 12, 12, 23, 23, 0, 1, time.Local),
				CompletionTime: time.Date(2020, 12, 12, 23, 24, 0, 1, time.Local),
				OriginalID:     12,
				ProcessedID:    13,
				Ratio:          0.5,
				OriginalType:   "jpeg",
				ProcessedType:  "png",
			},
			wantErr: nil,
		},
		{
			testName: "no such row in db",
			userID:   12,
			reqID:    19,
			repoReq:  nil,
			wantReq:  nil,
			wantErr:  sql.ErrNoRows,
		},
	}

	query := fmt.Sprintf(`SELECT id, op_status, request_time, completion_time, original_id,
	 processed_id, ratio, original_type, processed_type FROM %s WHERE id = .+ and user_id = .+`, repository.RequestTable)

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			repo := repository.NewReqPostgres(db)

			rows := sqlmock.NewRows([]string{"id", "op_status", "request_time", "completion_time", "original_id", "processed_id",
				"ratio", "original_type", "processed_type"})

			if tc.repoReq != nil {
				rows = rows.AddRow(tc.repoReq.ID, tc.repoReq.OpStatus, tc.repoReq.RequestTime, tc.repoReq.CompletionTime,
					tc.repoReq.OriginalID, tc.repoReq.ProcessedID, tc.repoReq.Ratio,
					tc.repoReq.OriginalType, tc.repoReq.ProcessedType)
			}

			mock.ExpectQuery(query).WithArgs(tc.reqID, tc.userID).WillReturnRows(rows)

			gotRequest, gotErr := repo.GetRequest(context.Background(), tc.userID, tc.reqID)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("Want error : %v, got error: %v", tc.wantErr, gotErr)
			}
			if !cmp.Equal(gotRequest, tc.wantReq) {
				t.Errorf("Want request: %v, got request: %v", tc.wantReq, gotRequest)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %s", err)
			}
		})
	}
}

func TestAddRequest(t *testing.T) {
	testCases := []struct {
		testName string
		userID   int
		req      *model.Request
		repoID   int
		wantID   int
		wantErr  error
	}{
		{
			testName: "all is good",
			userID:   12,
			req: &model.Request{
				ID:             24,
				OpStatus:       "done",
				RequestTime:    time.Date(2020, 12, 12, 23, 23, 0, 1, time.Local),
				CompletionTime: time.Date(2020, 12, 12, 23, 24, 0, 1, time.Local),
				OriginalID:     12,
				ProcessedID:    13,
				Ratio:          0.5,
				OriginalType:   "jpeg",
				ProcessedType:  "png",
			},
			repoID:  23,
			wantID:  23,
			wantErr: nil,
		},
	}

	query := fmt.Sprintf(`INSERT INTO %s \(op_status, request_time, original_id, 
		user_id, ratio, original_type, processed_type\)
		VALUES (.+, .+, .+, .+, .+, .+, .+) RETURNING id;`, repository.RequestTable)

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			repo := repository.NewReqPostgres(db)

			rows := sqlmock.NewRows([]string{"id"})

			if tc.req != nil {
				rows = rows.AddRow(tc.repoID)
			}

			mock.ExpectQuery(query).WithArgs(tc.req.OpStatus, tc.req.RequestTime,
				tc.req.OriginalID, tc.userID, tc.req.Ratio, tc.req.OriginalType,
				tc.req.ProcessedType).WillReturnRows(rows)

			gotID, gotErr := repo.AddRequest(context.Background(), tc.req, tc.userID)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("Want error : %v, got error: %v", tc.wantErr, gotErr)
			}
			if gotID != tc.wantID {
				t.Errorf("Want id: %v, got id: %v", tc.wantID, gotID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %s", err)
			}
		})
	}
}

func TestUpdateRequestStatus(t *testing.T) {
	testCases := []struct {
		testName string
		reqID    int
		status   string
		repoID   int
		wantErr  error
	}{
		{
			testName: "all is good",
			reqID:    12,
			status:   "done",
			repoID:   23,
			wantErr:  nil,
		},
		{
			testName: "such row not present in database",
			reqID:    12,
			status:   "done",
			repoID:   0,
			wantErr:  sql.ErrNoRows,
		},
	}

	query := fmt.Sprintf(`UPDATE %s SET op_status = .+ WHERE id = .+ RETURNING id`, repository.RequestTable)

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			repo := repository.NewReqPostgres(db)

			rows := sqlmock.NewRows([]string{"id"})

			if tc.repoID != 0 {
				rows = rows.AddRow(tc.repoID)
			}

			mock.ExpectQuery(query).WithArgs(tc.status, tc.reqID).WillReturnRows(rows)

			gotErr := repo.UpdateRequestStatus(context.Background(), tc.reqID, tc.status)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("Want error : %v, got error: %v", tc.wantErr, gotErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %s", err)
			}
		})
	}
}

func TestAddProcessedImageIDToRequest(t *testing.T) {
	testCases := []struct {
		testName string
		reqID    int
		imageID  int
		repoID   int
		wantErr  error
	}{
		{
			testName: "all is good",
			reqID:    12,
			imageID:  13,
			repoID:   23,
			wantErr:  nil,
		},
		{
			testName: "such row not present in database",
			reqID:    12,
			imageID:  13,
			repoID:   0,
			wantErr:  sql.ErrNoRows,
		},
	}

	query := fmt.Sprintf(`UPDATE %s SET processed_id = .+ WHERE id = .+ RETURNING id;`, repository.RequestTable)

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			repo := repository.NewReqPostgres(db)

			rows := sqlmock.NewRows([]string{"id"})

			if tc.repoID != 0 {
				rows = rows.AddRow(tc.repoID)
			}

			mock.ExpectQuery(query).WithArgs(tc.imageID, tc.reqID).WillReturnRows(rows)

			gotErr := repo.AddProcessedImageIDToRequest(context.Background(), tc.reqID, tc.imageID)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("Want error : %v, got error: %v", tc.wantErr, gotErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %s", err)
			}
		})
	}
}

func TestAddProcessedTimeToRequest(t *testing.T) {
	testCases := []struct {
		testName string
		reqID    int
		procTime time.Time
		repoID   int
		wantErr  error
	}{
		{
			testName: "all is good",
			reqID:    12,
			procTime: time.Date(2012, 3, 12, 3, 23, 3, 4, time.Local),
			repoID:   23,
			wantErr:  nil,
		},
		{
			testName: "such row not present in database",
			reqID:    12,
			procTime: time.Date(2012, 3, 12, 3, 23, 3, 4, time.Local),
			repoID:   0,
			wantErr:  sql.ErrNoRows,
		},
	}

	query := fmt.Sprintf(`UPDATE %s SET completion_time = .+ WHERE id = .+ RETURNING id;`, repository.RequestTable)

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			repo := repository.NewReqPostgres(db)

			rows := sqlmock.NewRows([]string{"id"})

			if tc.repoID != 0 {
				rows = rows.AddRow(tc.repoID)
			}

			mock.ExpectQuery(query).WithArgs(tc.procTime, tc.reqID).WillReturnRows(rows)

			gotErr := repo.AddProcessedTimeToRequest(context.Background(), tc.reqID, tc.procTime)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("Want error : %v, got error: %v", tc.wantErr, gotErr)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %s", err)
			}
		})
	}
}

func TestAddImage(t *testing.T) {
	testCases := []struct {
		testName  string
		userID    int
		imageInfo *model.Info
		repoID    int
		wantID    int
		wantErr   error
	}{
		{
			testName: "all is good",
			userID:   12,
			imageInfo: &model.Info{
				ResoultionX: 2040,
				ResoultionY: 1020,
				Type:        "png",
				URL:         "url to image",
			},
			repoID:  23,
			wantID:  23,
			wantErr: nil,
		},
	}

	query := fmt.Sprintf(`INSERT INTO %s \(resoolution_x, resoolution_y, im_type, image_url, user_id\)
		VALUES (.+, .+, .+, .+, .+) RETURNING id;`, repository.ImageTable)

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			repo := repository.NewReqPostgres(db)

			rows := sqlmock.NewRows([]string{"id"})

			if tc.imageInfo != nil {
				rows = rows.AddRow(tc.repoID)
			}

			mock.ExpectQuery(query).WithArgs(tc.imageInfo.ResoultionX, tc.imageInfo.ResoultionY,
				tc.imageInfo.Type, tc.imageInfo.URL, tc.userID).WillReturnRows(rows)

			gotID, gotErr := repo.AddImage(context.Background(), tc.userID, *tc.imageInfo)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("Want error : %v, got error: %v", tc.wantErr, gotErr)
			}
			if gotID != tc.wantID {
				t.Errorf("Want id: %v, got id: %v", tc.wantID, gotID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %s", err)
			}
		})
	}
}
