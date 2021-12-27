package service_test

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/Dyleme/image-coverter/internal/service"
	"github.com/Dyleme/image-coverter/internal/service/mocks"
	"github.com/golang/mock/gomock"
)

type genHashMock struct{}

func (g *genHashMock) GeneratePasswordHash(password string) string {
	return password
}

func (g *genHashMock) IsValidPassword(password string, hash []byte) bool {
	return password == string(hash)
}

type genJwtMock struct{}

var valToGetError = 256

var errCreateToken = errors.New("create token error")

func (g *genJwtMock) CreateToken(ctx context.Context, id int) (string, error) {
	if id == valToGetError {
		return "", errCreateToken
	}

	return "jwt" + strconv.Itoa(id), nil
}

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		testName  string
		user      model.User
		repoID    int
		repoError error
		wantID    int
		wantError error
	}{
		{
			testName: "All is good",
			user: model.User{
				Nickname: "Alike",
				Password: "hello",
			},
			repoID:    23,
			repoError: nil,
			wantID:    23,
			wantError: nil,
		},
		{
			testName: "Repository returned error",
			user: model.User{
				Nickname: "Alike",
				Password: "hello",
			},
			repoID:    0,
			repoError: errRepository,
			wantID:    0,
			wantError: errRepository,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()
			mockAuth := mocks.NewMockAutharizater(mockCtr)

			generator := &genHashMock{}
			jwtGnrt := &genJwtMock{}
			srvc := service.NewAuthSevice(mockAuth, generator, jwtGnrt)

			ctx := context.Background()

			mockAuth.EXPECT().CreateUser(ctx, tc.user).Return(tc.repoID, tc.repoError).Times(1)

			gotID, gotErr := srvc.CreateUser(ctx, tc.user)

			if !errors.Is(gotErr, tc.wantError) {
				t.Errorf("want error %b, got errror %v", tc.wantError, gotErr)
			}

			if gotID != tc.wantID {
				t.Errorf("want id %v, got id %v", tc.wantID, gotID)
			}
		})
	}
}

func TestValidateUser(t *testing.T) {
	testCases := []struct {
		testName   string
		user       model.User
		repoID     int
		repoPsswrd []byte
		repoError  error
		wantJwt    string
		wantError  error
	}{
		{
			testName: "All is good",
			user: model.User{
				Nickname: "Alike",
				Password: "123",
			},
			repoID:     25,
			repoPsswrd: []byte("123"),
			repoError:  nil,
			wantJwt:    "jwt" + strconv.Itoa(25),
			wantError:  nil,
		},
		{
			testName: "Repository error",
			user: model.User{
				Nickname: "Alike",
				Password: "123",
			},
			repoID:     0,
			repoPsswrd: nil,
			repoError:  errRepository,
			wantJwt:    "",
			wantError:  errRepository,
		},
		{
			testName: "Wrong password",
			user: model.User{
				Nickname: "Alike",
				Password: "123",
			},
			repoID:     19,
			repoPsswrd: []byte("321"),
			repoError:  nil,
			wantJwt:    "",
			wantError:  service.ErrWrongPassword,
		},
		{
			testName: "Token creation error",
			user: model.User{
				Nickname: "Alike",
				Password: "123",
			},
			repoID:     valToGetError,
			repoPsswrd: []byte("123"),
			repoError:  nil,
			wantJwt:    "",
			wantError:  errCreateToken,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			mockCtr := gomock.NewController(t)
			defer mockCtr.Finish()

			mockAuth := mocks.NewMockAutharizater(mockCtr)

			generator := &genHashMock{}
			jwtGnrt := &genJwtMock{}

			srvc := service.NewAuthSevice(mockAuth, generator, jwtGnrt)

			ctx := context.Background()
			mockAuth.EXPECT().GetPasswordHashAndID(ctx, tc.user.Nickname).Return(tc.repoPsswrd, tc.repoID, tc.repoError).Times(1)

			gotJwt, gotErr := srvc.ValidateUser(ctx, tc.user)

			if !errors.Is(gotErr, tc.wantError) {
				t.Errorf("want error %v, got error %v", tc.wantError, gotErr)
			}

			if gotJwt != tc.wantJwt {
				t.Errorf("want jwt %v, got jwt %v", tc.wantJwt, gotJwt)
			}
		})
	}
}
