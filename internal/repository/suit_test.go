package repository_test

import (
	"database/sql"

	"github.com/Dyleme/image-coverter/internal/repository"
	migrate "github.com/golang-migrate/migrate/v4"
	_pq "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var testDatabaseName = "postgres"

type PostgresSuit struct {
	suite.Suite
	DBConn          *sql.DB
	Migration       *migrate.Migrate
	MigrationFolder string
	conf            repository.DBConfig
	DSN             string
}

func (p *PostgresSuit) SetupSuite() {
	var err error

	p.DBConn, err = repository.NewPostgresDB(&p.conf)
	if err != nil {
		panic(err)
	}

	p.Migration, err = runMigration(p.DBConn,
		"./migration")
	require.NoError(p.T(), err)
}

func (p *PostgresSuit) TearDownSuite() {
	p.DBConn.Close()
}

func runMigration(dbConn *sql.DB, migrationFolder string) (*migrate.Migrate, error) {
	dataPath := "file://" + migrationFolder

	driver, err := _pq.WithInstance(dbConn, &_pq.Config{})
	if err != nil {
		return nil, err
	}

	migration, err := migrate.NewWithDatabaseInstance(dataPath, testDatabaseName, driver)
	if err != nil {
		return nil, err
	}

	return migration, nil
}
