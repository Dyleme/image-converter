package repository

import (
	"context"
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

const (
	UsersTable   = "users"
	RequestTable = "requests"
	ImageTable   = "images"
)

const (
	StatusQueued     = `queued`
	StatusProcessing = `processing`
	StatusDone       = `done`
)

// Config to connect to the database.
type DBConfig struct {
	UserName string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

type TxDB struct {
	*sql.DB
}

// Constuctor to the postgres database.
func NewPostgresDB(conf *DBConfig) (*sql.DB, error) {
	var db *sql.DB

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		conf.Host, conf.Port, conf.UserName, conf.Password, conf.DBName, conf.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Println(err)
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func oneRowInResult(result sql.Result) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	if rows != 1 {
		return fmt.Errorf("repo: %w", &NotSingleRowAffectedError{int(rows)})
	}

	return nil
}

func (db *TxDB) inTx(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		if err1 := tx.Rollback(); err1 != nil {
			return fmt.Errorf("rolling back transaction %v, (original error %v)",
				err1, err) //nolint:errorlint // making combined error
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}
