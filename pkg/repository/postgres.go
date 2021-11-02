package repository

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	usersTable   = "users"
	requestTable = "requests"
	imageTable   = "images"
)

type DBConfig struct {
	UserName string
	Password string
	Host     string
	Port     string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(conf *DBConfig) (*sql.DB, error) {
	var db *sql.DB

	connStr := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=%v",
		conf.UserName, conf.Password, conf.DBName, conf.SSLMode)

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
