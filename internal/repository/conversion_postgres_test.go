package repository_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/Dyleme/image-coverter/internal/repository"
	"github.com/stretchr/testify/assert"
)

func NewConvMock(t *testing.T) (*repository.ConvPostgres, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	repo := repository.NewConvPostgres(db)

	return repo, mock
}

var addImageWithResolutionQuery = regexp.QuoteMeta(fmt.Sprintf(`INSERT INTO %s 
(im_type, image_url, user_id, resoolution_x, resoolution_y)
VALUES ($1, $2, $3, $4, $5) RETURNING id`, repository.ImageTable))
var updateRequestStatusQuery = fmt.Sprintf(`UPDATE %s SET op_status = .+ 
WHERE id = .+`, repository.RequestTable)
var addProcessedIDQuery = fmt.Sprintf(`UPDATE %s SET processed_id = .+ 
WHERE id = .+`, repository.RequestTable)
var addProcessedTimeQuery = fmt.Sprintf(`UPDATE %s SET completion_time = .+ 
WHERE id = .+`, repository.RequestTable)

var (
	errAddImageToDB     = errors.New("repo add image to db error")
	errAddProcessedID   = errors.New("repo add processed id")
	errAddProcessedTime = errors.New("repo add processed time error")
	errUpdateStatus     = errors.New("repo update status error")
)

func TestConvPostgres_AddImageDB(t *testing.T) {
	testCases := []struct {
		testName string
		userID   int
		reqID    int
		imgInfo  *model.ReuquestImageInfo
		status   string
		width    int
		height   int
		time     time.Time
		initMock func(sqlmock.Sqlmock, int, int, *model.ReuquestImageInfo, int, int, string, time.Time) sqlmock.Sqlmock
		wantErr  error
	}{
		{
			testName: "all is good",
			userID:   2,
			reqID:    3,
			imgInfo: &model.ReuquestImageInfo{
				Type: "jpeg",
				URL:  "image url",
			},
			status: repository.StatusDone,
			time:   time.Date(2021, 1, 4, 10, 25, 34, 0, &time.Location{}),
			initMock: func(mock sqlmock.Sqlmock, user, req int, imgInfo *model.ReuquestImageInfo,
				width, height int, status string, t time.Time) sqlmock.Sqlmock {
				imageID := 32
				imageRow := RepoReturnID(imageID)
				mock.ExpectBegin()
				mock.ExpectQuery(addImageWithResolutionQuery).WithArgs(imgInfo.Type, imgInfo.URL,
					user, width, height).WillReturnRows(imageRow)
				mock.ExpectExec(addProcessedIDQuery).WithArgs(imageID, req).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(addProcessedTimeQuery).WithArgs(t, req).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(updateRequestStatusQuery).WithArgs(repository.StatusDone, req).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
				return mock
			},
			wantErr: nil,
		},
		{
			testName: "error at update status",
			userID:   2,
			reqID:    3,
			imgInfo: &model.ReuquestImageInfo{
				Type: "jpeg",
				URL:  "image url",
			},
			status: repository.StatusDone,
			time:   time.Date(2021, 1, 4, 10, 25, 34, 0, &time.Location{}),
			initMock: func(mock sqlmock.Sqlmock, user, req int, imgInfo *model.ReuquestImageInfo,
				width, height int, status string, t time.Time) sqlmock.Sqlmock {
				imageID := 32
				imageRow := RepoReturnID(imageID)
				mock.ExpectBegin()
				mock.ExpectQuery(addImageWithResolutionQuery).WithArgs(imgInfo.Type, imgInfo.URL,
					user, width, height).WillReturnRows(imageRow)
				mock.ExpectExec(addProcessedIDQuery).WithArgs(imageID, req).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(addProcessedTimeQuery).WithArgs(t, req).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(updateRequestStatusQuery).WithArgs(repository.StatusDone, req).
					WillReturnResult(sqlmock.NewErrorResult(errUpdateStatus))
				mock.ExpectRollback()
				return mock
			},
			wantErr: errUpdateStatus,
		},
		{
			testName: "error at add time",
			userID:   2,
			reqID:    3,
			imgInfo: &model.ReuquestImageInfo{
				Type: "jpeg",
				URL:  "image url",
			},
			status: repository.StatusDone,
			time:   time.Date(2021, 1, 4, 10, 25, 34, 0, &time.Location{}),
			initMock: func(mock sqlmock.Sqlmock, user, req int, imgInfo *model.ReuquestImageInfo,
				width, height int, status string, t time.Time) sqlmock.Sqlmock {
				imageID := 32
				imageRow := RepoReturnID(imageID)
				mock.ExpectBegin()
				mock.ExpectQuery(addImageWithResolutionQuery).WithArgs(imgInfo.Type, imgInfo.URL,
					user, width, height).WillReturnRows(imageRow)
				mock.ExpectExec(addProcessedIDQuery).WithArgs(imageID, req).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(addProcessedTimeQuery).WithArgs(t, req).
					WillReturnResult(sqlmock.NewErrorResult(errAddProcessedTime))
				mock.ExpectRollback()
				return mock
			},
			wantErr: errAddProcessedTime,
		},
		{
			testName: "error at add id",
			userID:   2,
			reqID:    3,
			imgInfo: &model.ReuquestImageInfo{
				Type: "jpeg",
				URL:  "image url",
			},
			status: repository.StatusDone,
			time:   time.Date(2021, 1, 4, 10, 25, 34, 0, &time.Location{}),
			initMock: func(mock sqlmock.Sqlmock, user, req int, imgInfo *model.ReuquestImageInfo,
				width, height int, status string, t time.Time) sqlmock.Sqlmock {
				imageID := 32
				imageRow := RepoReturnID(imageID)
				mock.ExpectBegin()
				mock.ExpectQuery(addImageWithResolutionQuery).WithArgs(imgInfo.Type, imgInfo.URL,
					user, width, height).WillReturnRows(imageRow)
				mock.ExpectExec(addProcessedIDQuery).WithArgs(imageID, req).
					WillReturnResult(sqlmock.NewErrorResult(errAddProcessedID))
				mock.ExpectRollback()
				return mock
			},
			wantErr: errAddProcessedID,
		},
		{
			testName: "error at add image to requests",
			userID:   2,
			reqID:    3,
			imgInfo: &model.ReuquestImageInfo{
				Type: "jpeg",
				URL:  "image url",
			},
			status: repository.StatusDone,
			time:   time.Date(2021, 1, 4, 10, 25, 34, 0, &time.Location{}),
			initMock: func(mock sqlmock.Sqlmock, user, req int, imgInfo *model.ReuquestImageInfo,
				width, height int, status string, t time.Time) sqlmock.Sqlmock {
				mock.ExpectBegin()
				mock.ExpectQuery(addImageWithResolutionQuery).WithArgs(imgInfo.Type, imgInfo.URL,
					user, width, height).WillReturnError(errAddImageToDB)
				mock.ExpectRollback()
				return mock
			},
			wantErr: errAddImageToDB,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			repo, mock := NewConvMock(t)

			mock = tc.initMock(mock, tc.userID, tc.reqID,
				tc.imgInfo, tc.width, tc.height, tc.status, tc.time)

			err := repo.AddProcessedImage(context.Background(), tc.userID, tc.reqID, tc.imgInfo,
				tc.width, tc.height, tc.status, tc.time)

			assert.ErrorIs(t, err, tc.wantErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %v", err)
			}
		})
	}
}
