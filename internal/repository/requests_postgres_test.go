package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/Dyleme/image-coverter/internal/repository"
	"github.com/stretchr/testify/assert"
)

func NewReqMock(t *testing.T) (*repository.ReqPostgres, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	repo := repository.NewReqPostgres(db)

	return repo, mock
}

func RepoReturnID(repo *repository.ReqPostgres, id int) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id"})
	if id != 0 {
		rows = rows.AddRow(id)
	}

	return rows
}

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
			repo, mock := NewReqMock(t)

			rows := sqlmock.NewRows([]string{"id", "op_status", "request_time", "completion_time", "original_id", "processed_id",
				"ratio", "original_type", "processed_type"})

			if tc.repoReq != nil {
				rows = rows.AddRow(tc.repoReq.ID, tc.repoReq.OpStatus, tc.repoReq.RequestTime, tc.repoReq.CompletionTime,
					tc.repoReq.OriginalID, tc.repoReq.ProcessedID, tc.repoReq.Ratio,
					tc.repoReq.OriginalType, tc.repoReq.ProcessedType)
			}

			mock.ExpectQuery(query).WithArgs(tc.reqID, tc.userID).WillReturnRows(rows)

			gotRequest, gotErr := repo.GetRequest(context.Background(), tc.userID, tc.reqID)

			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, gotRequest, tc.wantReq)

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
			repo, mock := NewReqMock(t)

			rows := RepoReturnID(repo, tc.repoID)

			mock.ExpectQuery(query).WithArgs(tc.req.OpStatus, tc.req.RequestTime,
				tc.req.OriginalID, tc.userID, tc.req.Ratio, tc.req.OriginalType,
				tc.req.ProcessedType).WillReturnRows(rows)

			gotID, gotErr := repo.AddRequest(context.Background(), tc.req, tc.userID)

			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, gotID, tc.wantID)

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
		imageInfo *model.ReuquestImageInfo
		repoID    int
		wantID    int
		wantErr   error
	}{
		{
			testName: "all is good",
			userID:   12,
			imageInfo: &model.ReuquestImageInfo{
				Type: "png",
				URL:  "url to image",
			},
			repoID:  23,
			wantID:  23,
			wantErr: nil,
		},
	}

	query := fmt.Sprintf(`INSERT INTO %s \(im_type, image_url, user_id\)
		VALUES (.+, .+, .+) RETURNING id;`, repository.ImageTable)

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			repo, mock := NewReqMock(t)

			rows := RepoReturnID(repo, tc.repoID)

			mock.ExpectQuery(query).WithArgs(tc.imageInfo.Type, tc.imageInfo.URL, tc.userID).WillReturnRows(rows)

			gotID, gotErr := repo.AddImage(context.Background(), tc.userID, *tc.imageInfo)

			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, gotID, tc.wantID)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteRequest(t *testing.T) {
	testCases := []struct {
		testName  string
		userID    int
		reqID     int
		im1repoID int
		im2repoID int
		wantIm1ID int
		wantIm2ID int
		wantErr   error
	}{
		{
			testName:  "all is good",
			userID:    12,
			reqID:     13,
			im1repoID: 23,
			im2repoID: 24,
			wantIm1ID: 23,
			wantIm2ID: 24,
			wantErr:   nil,
		}, {
			testName:  "such rown not exist",
			userID:    12,
			reqID:     13,
			im1repoID: 0,
			im2repoID: 0,
			wantIm1ID: 0,
			wantIm2ID: 0,
			wantErr:   sql.ErrNoRows,
		}, {
			testName:  "only one id is exist in row",
			userID:    12,
			reqID:     13,
			im1repoID: 23,
			im2repoID: 0,
			wantIm1ID: 23,
			wantIm2ID: 0,
			wantErr:   nil,
		},
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = .+ AND id = .+ 
		RETURNING original_id, processed_id`, repository.RequestTable)

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			repo, mock := NewReqMock(t)

			rows := sqlmock.NewRows([]string{"original_id", "processed_id"})

			if tc.im1repoID != 0 {
				rows = rows.AddRow(tc.im1repoID, tc.im2repoID)
			}

			mock.ExpectQuery(query).WithArgs(tc.userID, tc.reqID).WillReturnRows(rows)

			gotIm1ID, gotIm2ID, gotErr := repo.DeleteRequest(context.Background(), tc.userID, tc.reqID)

			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, gotIm2ID, tc.wantIm2ID)
			assert.Equal(t, gotIm1ID, tc.wantIm1ID)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %s", err)
			}
		})
	}
}

func TestDeleteImage(t *testing.T) {
	testCases := []struct {
		testName string
		userID   int
		imageID  int
		repoURL  string
		wantURL  string
		wantErr  error
	}{
		{
			testName: "all is good",
			userID:   12,
			imageID:  23,
			repoURL:  "url to image",
			wantURL:  "url to image",
			wantErr:  nil,
		},
	}

	query := fmt.Sprintf(`DELETE FROM %s WHERE user_id = .+ AND id = .+ RETURNING image_url`, repository.ImageTable)

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			repo, mock := NewReqMock(t)

			rows := sqlmock.NewRows([]string{"image_url"})

			if tc.repoURL != "" {
				rows = rows.AddRow(tc.repoURL)
			}

			mock.ExpectQuery(query).WithArgs(tc.userID, tc.imageID).WillReturnRows(rows)

			gotURL, gotErr := repo.DeleteImage(context.Background(), tc.userID, tc.imageID)

			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, gotURL, tc.wantURL)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %s", err)
			}
		})
	}
}
