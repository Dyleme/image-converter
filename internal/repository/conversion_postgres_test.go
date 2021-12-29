package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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
			repo, mock := NewConvMock(t)

			rows := sqlmock.NewRows([]string{"id"})
			if tc.repoID != 0 {
				rows = rows.AddRow(tc.repoID)
			}

			mock.ExpectQuery(query).WithArgs(tc.status, tc.reqID).WillReturnRows(rows)

			gotErr := repo.UpdateRequestStatus(context.Background(), tc.reqID, tc.status)

			assert.ErrorIs(t, gotErr, tc.wantErr)

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
			reqID:    521,
			imageID:  13,
			repoID:   23,
			wantErr:  nil,
		},
		{
			testName: "such row not present in database",
			reqID:    932,
			imageID:  13,
			repoID:   0,
			wantErr:  sql.ErrNoRows,
		},
	}

	query := fmt.Sprintf(`UPDATE %s SET processed_id = .+ WHERE id = .+ RETURNING id;`, repository.RequestTable)

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			repo, mock := NewConvMock(t)

			rows := sqlmock.NewRows([]string{"id"})
			if tc.repoID != 0 {
				rows = rows.AddRow(tc.repoID)
			}

			mock.ExpectQuery(query).WithArgs(tc.imageID, tc.reqID).WillReturnRows(rows)

			gotErr := repo.AddProcessedImageIDToRequest(context.Background(), tc.reqID, tc.imageID)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %s", err)
			}

			assert.ErrorIs(t, gotErr, tc.wantErr)
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
			repo, mock := NewConvMock(t)

			rows := sqlmock.NewRows([]string{"id"})
			if tc.repoID != 0 {
				rows = rows.AddRow(tc.repoID)
			}

			mock.ExpectQuery(query).WithArgs(tc.procTime, tc.reqID).WillReturnRows(rows)

			gotErr := repo.AddProcessedTimeToRequest(context.Background(), tc.reqID, tc.procTime)

			assert.ErrorIs(t, gotErr, tc.wantErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %s", err)
			}
		})
	}
}
