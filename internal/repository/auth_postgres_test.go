package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/Dyleme/image-coverter/internal/repository"
	"github.com/stretchr/testify/assert"
)

func NewAuthMock(t *testing.T) (*repository.AuthPostgres, sqlmock.Sqlmock) {
	t.Helper()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	repo := repository.NewAuthPostgres(db)

	return repo, mock
}

var createUserQuery = fmt.Sprintf(`INSERT INTO %s (nickname, email, password_hash)
	 VALUES ($1, $2, $3) RETURNING id`, repository.UsersTable)

var errAlreadyExists = errors.New("such row already exists")

func TestAuthPostgres_CreateUser(t *testing.T) {
	testCases := []struct {
		testName string
		user     model.User
		initMock func(sqlmock.Sqlmock, *model.User) sqlmock.Sqlmock
		wantID   int
		wantErr  error
	}{
		{
			testName: "all is good",
			user: model.User{
				Nickname: "alekse",
				Email:    "my@email.com",
				Password: "password",
			},
			initMock: func(mock sqlmock.Sqlmock, user *model.User) sqlmock.Sqlmock {
				rows := sqlmock.NewRows([]string{"id"}).AddRow(12)
				mock.ExpectQuery(regexp.QuoteMeta(createUserQuery)).WithArgs(user.Nickname, user.Email, user.Password).
					WillReturnRows(rows)

				return mock
			},
			wantID:  12,
			wantErr: nil,
		},
		{
			testName: "such row already exists",
			user: model.User{
				Nickname: "alekse",
				Email:    "my@email.com",
				Password: "password",
			},
			initMock: func(mock sqlmock.Sqlmock, user *model.User) sqlmock.Sqlmock {
				mock.ExpectQuery(regexp.QuoteMeta(createUserQuery)).WithArgs(user.Nickname, user.Email, user.Password).
					WillReturnError(errAlreadyExists)

				return mock
			},
			wantID:  0,
			wantErr: errAlreadyExists,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			repo, mock := NewAuthMock(t)

			mock = tc.initMock(mock, &tc.user)

			gotID, gotErr := repo.CreateUser(context.Background(), tc.user)

			assert.ErrorIs(t, gotErr, tc.wantErr)
			assert.Equal(t, gotID, tc.wantID)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %s", err)
			}
		})
	}
}

var getPasswrodAndIDQuery = fmt.Sprintf("SELECT password_hash, id FROM %s WHERE nickname = ?", repository.UsersTable)

func TestAuthPostgres_GetPasswordAndID(t *testing.T) {
	testCases := []struct {
		testName     string
		userNickname string
		initMock     func(sqlmock.Sqlmock, string) sqlmock.Sqlmock
		wantID       int
		wantPassword []byte
		wantErr      error
	}{
		{
			testName:     "all is good",
			userNickname: "alekse",
			initMock: func(mock sqlmock.Sqlmock, nickname string) sqlmock.Sqlmock {
				rows := sqlmock.NewRows([]string{"password_hash", "id"})
				rows.AddRow("password", 12)

				mock.ExpectQuery(getPasswrodAndIDQuery).WithArgs(nickname).
					WillReturnRows(rows)

				return mock
			},
			wantID:       12,
			wantPassword: []byte("password"),
			wantErr:      nil,
		},
		{
			testName:     "unknown nickname",
			userNickname: "unknown",
			initMock: func(mock sqlmock.Sqlmock, nickname string) sqlmock.Sqlmock {
				mock.ExpectQuery(getPasswrodAndIDQuery).WithArgs(nickname).
					WillReturnError(sql.ErrNoRows)

				return mock
			},
			wantID:       0,
			wantPassword: nil,
			wantErr:      sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			repo, mock := NewAuthMock(t)

			mock = tc.initMock(mock, tc.userNickname)

			gotPassword, gotID, gotErr := repo.GetPasswordHashAndID(context.Background(), tc.userNickname)

			assert.Equal(t, gotID, tc.wantID)
			assert.Equal(t, gotPassword, tc.wantPassword)
			assert.ErrorIs(t, gotErr, tc.wantErr)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were fulfilled expectations: %s", err)
			}
		})
	}
}
