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

// TxDb is a struct which compose *sql.DB.
// It is made to provide inTx(...) method  to make querys in transaction.
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

// NotSingleRowAffectedError is error which is used when not one raw came from Query.
type NotSingleRowAffectedError struct {
	amountAffected int
}

func (e *NotSingleRowAffectedError) Error() string {
	return fmt.Sprintf("expected single row affected, got %v rows affected", e.amountAffected)
}

// oneRowInResult is a function check if result has one row,
// returns NotSingleRowAffectedError if no.
func oneRowInResult(result sql.Result) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows != 1 {
		return &NotSingleRowAffectedError{int(rows)}
	}

	return nil
}

// inTx is method which allows you to make queries in transaction.
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
