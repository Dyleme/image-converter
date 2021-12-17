package repository

import (
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
