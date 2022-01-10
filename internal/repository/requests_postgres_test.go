package repository_test

import (
	"context"
	"database/sql"
	"errors"
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

func RepoReturnID(id int) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{"id"})
	if id != 0 {
		rows = rows.AddRow(id)
	}

	return rows
}

var getRequestQuery = fmt.Sprintf(`SELECT id, op_status, request_time, completion_time, original_id,
	 processed_id, ratio, original_type, processed_type FROM %s WHERE id = .+ and user_id = .+`, repository.RequestTable)

func TestReqPostgres_GetRequest(t *testing.T) {
	testCases := []struct {
		testName string
		userID   int
		reqID    int
		initMock func(sqlmock.Sqlmock, int, int, *model.Request) sqlmock.Sqlmock
		wantReq  *model.Request
		wantErr  error
	}{
		{
			testName: "all is good",
			userID:   12,
			reqID:    19,
			initMock: func(mock sqlmock.Sqlmock, userID, reqID int, req *model.Request) sqlmock.Sqlmock {
				rows := sqlmock.NewRows([]string{"id", "op_status", "request_time", "completion_time",
					"original_id", "processed_id", "ratio", "original_type", "processed_type"})

				rows = rows.AddRow(req.ID, req.OpStatus, req.RequestTime, req.CompletionTime,
					req.OriginalID, req.ProcessedID, req.Ratio,
					req.OriginalType, req.ProcessedType)

				mock.ExpectQuery(getRequestQuery).WithArgs(reqID, userID).
					WillReturnRows(rows)
				return mock
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
			initMock: func(mock sqlmock.Sqlmock, userID, reqID int, req *model.Request) sqlmock.Sqlmock {
				mock.ExpectQuery(getRequestQuery).WithArgs(reqID, userID).
					WillReturnError(sql.ErrNoRows)
				return mock
			},
			wantReq: nil,
			wantErr: sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			repo, mock := NewReqMock(t)

			mock = tc.initMock(mock, tc.userID, tc.reqID, tc.wantReq)

			gotRequest, gotErr := repo.GetRequest(context.Background(), tc.userID, tc.reqID)

			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, gotRequest, tc.wantReq)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %s", err)
			}
		})
	}
}

var (
	addImageQuery = fmt.Sprintf(`INSERT INTO %s \(im_type, image_url, user_id\)
		VALUES (.+, .+, .+) RETURNING id;`, repository.ImageTable)

	addRequestQuery = fmt.Sprintf(`INSERT INTO %s \(op_status, request_time, original_id, 
		user_id, ratio, original_type, processed_type\)
		VALUES (.+, .+, .+, .+, .+, .+, .+) RETURNING id;`, repository.RequestTable)
)

var (
	errAddingImage   = errors.New("error while adding image")
	errAddingRequest = errors.New("error while adding reqeust")
)

func TestReqPostgres_AddImageAndRequest(t *testing.T) {
	testCases := []struct {
		testName  string
		userID    int
		reqID     int
		reqInfo   *model.Request
		imageInfo *model.ReuquestImageInfo
		initMock  func(int, *model.ReuquestImageInfo,
			*model.Request) (*repository.ReqPostgres, sqlmock.Sqlmock)
		imageID int
		wantID  int
		wantErr error
	}{
		{
			testName: "all is good",
			userID:   12,
			imageInfo: &model.ReuquestImageInfo{
				Type: "jpeg",
				URL:  "image url",
			},
			reqInfo: &model.Request{
				OpStatus:      repository.StatusDone,
				RequestTime:   time.Date(2022, 1, 3, 14, 36, 2, 32, &time.Location{}),
				OriginalID:    26,
				Ratio:         0.5,
				OriginalType:  "jpeg",
				ProcessedType: "type",
			},
			initMock: func(userID int, im *model.ReuquestImageInfo,
				req *model.Request) (*repository.ReqPostgres, sqlmock.Sqlmock) {
				repo, mock := NewReqMock(t)
				imageRow := RepoReturnID(req.OriginalID)
				reqRow := RepoReturnID(13)

				mock.ExpectBegin()

				mock.ExpectQuery(addImageQuery).WithArgs(im.Type, im.URL, userID).
					WillReturnRows(imageRow)
				mock.ExpectQuery(addRequestQuery).WithArgs(req.OpStatus, req.RequestTime,
					req.OriginalID, userID, req.Ratio,
					req.OriginalType, req.ProcessedType).
					WillReturnRows(reqRow)

				mock.ExpectCommit()

				return repo, mock
			},
			wantID:  13,
			wantErr: nil,
		},
		{
			testName: "transaction rollback image error",
			userID:   12,
			imageInfo: &model.ReuquestImageInfo{
				Type: "jpeg",
				URL:  "image url",
			},
			reqInfo: &model.Request{
				OpStatus:      repository.StatusDone,
				RequestTime:   time.Date(2022, 1, 3, 14, 36, 2, 32, &time.Location{}),
				OriginalID:    26,
				Ratio:         0.5,
				OriginalType:  "jpeg",
				ProcessedType: "type",
			},
			initMock: func(userID int, im *model.ReuquestImageInfo,
				req *model.Request) (*repository.ReqPostgres, sqlmock.Sqlmock) {
				repo, mock := NewReqMock(t)

				mock.ExpectBegin()

				mock.ExpectQuery(addImageQuery).WithArgs(im.Type, im.URL, userID).
					WillReturnError(errAddingImage)

				mock.ExpectRollback()

				return repo, mock
			},
			wantID:  0,
			wantErr: errAddingImage,
		},
		{
			testName: "transaction rollback request error",
			userID:   12,
			imageInfo: &model.ReuquestImageInfo{
				Type: "jpeg",
				URL:  "image url",
			},
			reqInfo: &model.Request{
				OpStatus:      repository.StatusDone,
				RequestTime:   time.Date(2022, 1, 3, 14, 36, 2, 32, &time.Location{}),
				OriginalID:    26,
				Ratio:         0.5,
				OriginalType:  "jpeg",
				ProcessedType: "type",
			},
			initMock: func(userID int, im *model.ReuquestImageInfo,
				req *model.Request) (*repository.ReqPostgres, sqlmock.Sqlmock) {
				repo, mock := NewReqMock(t)
				imageRow := RepoReturnID(req.OriginalID)

				mock.ExpectBegin()

				mock.ExpectQuery(addImageQuery).WithArgs(im.Type, im.URL, userID).
					WillReturnRows(imageRow)
				mock.ExpectQuery(addRequestQuery).WithArgs(req.OpStatus, req.RequestTime,
					req.OriginalID, userID, req.Ratio,
					req.OriginalType, req.ProcessedType).
					WillReturnError(errAddingRequest)

				mock.ExpectRollback()

				return repo, mock
			},
			wantID:  0,
			wantErr: errAddingRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			repo, mock := tc.initMock(tc.userID, tc.imageInfo, tc.reqInfo)

			reqID, gotErr := repo.AddImageAndRequest(context.Background(), tc.userID,
				tc.imageInfo, tc.reqInfo)

			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, reqID, tc.wantID)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %v", err)
			}
		})
	}
}

var delteRequestQuery = fmt.Sprintf(`DELETE FROM %s WHERE user_id = .+ AND id = .+ 
		RETURNING original_id, processed_id`, repository.RequestTable)

var deleteImageQuery = fmt.Sprintf(`DELETE FROM %s WHERE user_id = .+ AND id = .+
		RETURNING image_url`, repository.ImageTable)

func TestReqPostgres_DeleteRequestAndImage(t *testing.T) {
	testCases := []struct {
		testName   string
		userID     int
		reqID      int
		initMock   func(sqlmock.Sqlmock, int, int) sqlmock.Sqlmock
		wantIm1URL string
		wantIm2URL string
		wantErr    error
	}{
		{
			testName: "all is good",
			userID:   12,
			reqID:    13,
			initMock: func(mock sqlmock.Sqlmock, userID, reqID int) sqlmock.Sqlmock {
				idRows := sqlmock.NewRows([]string{"original_id", "processed_id"})
				idRows.AddRow(23, 24)

				url1Row := sqlmock.NewRows([]string{"image_url"}).AddRow("im 1 url")
				url2Row := sqlmock.NewRows([]string{"image_url"}).AddRow("im 2 url")
				mock.ExpectBegin()
				mock.ExpectQuery(delteRequestQuery).WithArgs(userID, reqID).
					WillReturnRows(idRows)
				mock.ExpectQuery(deleteImageQuery).WithArgs(userID, 23).
					WillReturnRows(url1Row)
				mock.ExpectQuery(deleteImageQuery).WithArgs(userID, 24).
					WillReturnRows(url2Row)
				mock.ExpectCommit()
				return mock
			},
			wantIm1URL: "im 1 url",
			wantIm2URL: "im 2 url",
			wantErr:    nil,
		},
		{
			testName: "such rown not exist",
			userID:   12,
			reqID:    13,
			initMock: func(mock sqlmock.Sqlmock, userID, reqID int) sqlmock.Sqlmock {
				mock.ExpectBegin()
				mock.ExpectQuery(delteRequestQuery).WithArgs(userID, reqID).
					WillReturnError(sql.ErrNoRows)
				mock.ExpectRollback()
				return mock
			},
			wantIm1URL: "",
			wantIm2URL: "",
			wantErr:    sql.ErrNoRows,
		},
		{
			testName: "only one id is exist in row",
			userID:   12,
			reqID:    13,
			initMock: func(mock sqlmock.Sqlmock, userID, reqID int) sqlmock.Sqlmock {
				idRows := sqlmock.NewRows([]string{"original_id", "processed_id"})
				idRows.AddRow(23, 0)

				url1Row := sqlmock.NewRows([]string{"image_url"}).AddRow("im 1 url")
				mock.ExpectBegin()
				mock.ExpectQuery(delteRequestQuery).WithArgs(userID, reqID).
					WillReturnRows(idRows)
				mock.ExpectQuery(deleteImageQuery).WithArgs(userID, 23).
					WillReturnRows(url1Row)
				mock.ExpectCommit()
				return mock
			},
			wantIm1URL: "im 1 url",
			wantIm2URL: "",
			wantErr:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			repo, mock := NewReqMock(t)

			mock = tc.initMock(mock, tc.userID, tc.reqID)

			gotIm1URL, gotIm2URL, gotErr := repo.DeleteRequestAndImage(context.Background(), tc.userID, tc.reqID)

			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, gotIm2URL, tc.wantIm2URL)
			assert.Equal(t, gotIm1URL, tc.wantIm1URL)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %s", err)
			}
		})
	}
}
