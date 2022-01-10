package repository_test

import (
	"context"
	"database/sql"
	"log"
	"testing"

	"github.com/Dyleme/image-coverter/internal/model"
	"github.com/Dyleme/image-coverter/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type PostgresAuthTest struct {
	PostgresSuit
}

func TestAuthSuit(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip category testing postgres")
	}

	categotySuit := &PostgresAuthTest{
		PostgresSuit{
			MigrationFolder: "migration",
			conf: repository.DBConfig{
				UserName: "postgres",
				Password: "postgres",
				Port:     "5432",
				Host:     "localhost",
				DBName:   testDatabaseName,
				SSLMode:  "disable",
			},
		},
	}

	suite.Run(t, categotySuit)
}

func (p *PostgresAuthTest) SetupTest() {
	log.Println("Starting a Test. Migrating the Database")

	err := p.Migration.Up()
	require.NoError(p.T(), err)
	log.Println("Database Migrated Successfully")
}

func (p *PostgresAuthTest) TearDownTest() {
	log.Println("Finifshing a Test. Dropping the Database")

	err := p.Migration.Down()
	require.NoError(p.T(), err)
	log.Println("Database Dropped Successfully")
}

func (p *PostgresAuthTest) TestGet() {
	repo := repository.NewAuthPostgres(p.DBConn)
	testCases := []struct {
		testName     string
		name         string
		wantPassword string
		wantID       int
		wantErr      error
	}{
		{
			testName:     "basic get",
			name:         "name2",
			wantPassword: "password2",
			wantID:       2,
			wantErr:      nil,
		},
		{
			testName:     "user not exist",
			name:         "not exist",
			wantPassword: "",
			wantID:       0,
			wantErr:      sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		p.T().Run(tc.testName, func(t *testing.T) {
			b, id, err := repo.GetPasswordHashAndID(context.TODO(), tc.name)
			assert.ErrorIs(p.T(), err, tc.wantErr)
			assert.Equal(p.T(), id, tc.wantID)
			if b != nil {
				assert.Equal(p.T(), string(b), tc.wantPassword)
			}
		})
	}
}

func (p *PostgresAuthTest) TestCreateUser() {
	repo := repository.NewAuthPostgres(p.DBConn)
	testCases := []struct {
		testName string
		user     model.User
		wantID   int
		wantErr  error
	}{
		{
			testName: "basic create",
			user: model.User{
				Nickname: "new name",
				Password: "password",
			},
			wantID:  3,
			wantErr: nil,
		},
		{
			testName: "nickname is used",
			user: model.User{
				Nickname: "name",
				Password: "password",
			},
			wantID:  0,
			wantErr: repository.ErrDuplicatedNickname,
		},
	}

	for _, tc := range testCases {
		p.T().Run(tc.testName, func(t *testing.T) {
			id, _ := repo.CreateUser(context.TODO(), tc.user)
			assert.Equal(p.T(), id, tc.wantID)
		})
	}
}
