package db

import (
	"context"
	"database/sql"
	"gosvelte/config"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func OpenDB(cfg config.Config) (*sql.DB, error) {

	// Open the connection with DSN
	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.DB.MaxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	log.Println("successfully connected to the database")

	return db, nil
}
