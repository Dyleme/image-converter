package repository_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Dyleme/image-coverter/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestDownloadPostgres_GetImageUrl(t *testing.T) {
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
			wantErr:  repository.ImageNotExistError{},
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

			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, gotURL, tc.wantURL)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %s", err)
			}
		})
	}
}
