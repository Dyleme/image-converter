package repository_test

import (
	"database/sql"
	"log"
	"time"

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

	timeout := 5 * time.Second
	timeoutExceeded := time.After(timeout)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutExceeded:
			p.T().Error("connecton timeout")
		case <-ticker.C:
			p.DBConn, err = repository.NewPostgresDB(&p.conf)
			if err != nil {
				log.Println("failed to connect to database")
			}

			p.Migration, err = runMigration(p.DBConn,
				"./migration")
			require.NoError(p.T(), err)

			return
		}
	}
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
