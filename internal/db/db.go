package db

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewDB(dsn string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxOpenConns)

	dur, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(time.Duration(dur))

	db.SetMaxIdleConns(maxIdleConns)

	return db, nil
}
