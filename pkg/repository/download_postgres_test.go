package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Dyleme/image-coverter/pkg/repository"
)

func TestGetImageUrl(t *testing.T) {
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
			imageID:  19,
			repoURL:  "url to image",
			wantURL:  "url to image",
			wantErr:  nil,
		},
		{
			testName: "no such row in db",
			userID:   12,
			imageID:  19,
			repoURL:  "",
			wantURL:  "",
			wantErr:  sql.ErrNoRows,
		},
	}

	query := fmt.Sprintf("SELECT image_url FROM %s WHERE user_id = .+ AND id = .+", repository.ImageTable)

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}

			repo := repository.NewDownloadPostgres(db)

			rows := sqlmock.NewRows([]string{"image_url"})
			if tc.repoURL != "" {
				rows = rows.AddRow(tc.repoURL)
			}

			mock.ExpectQuery(query).WithArgs(tc.userID, tc.imageID).WillReturnRows(rows)

			gotURL, gotErr := repo.GetImageURL(context.Background(), tc.userID, tc.imageID)

			if !errors.Is(gotErr, tc.wantErr) {
				t.Errorf("Want error : %v, got error: %v", tc.wantErr, gotErr)
			}

			if gotURL != tc.wantURL {
				t.Errorf("Want url : %v, got url: %v", tc.wantURL, gotURL)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %s", err)
			}
		})
	}
}
