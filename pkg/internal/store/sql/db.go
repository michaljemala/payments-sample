package sql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Config struct {
	Driver             string
	DSN                string
	ConnMaxLifetime    time.Duration
	MaxIdleConnections int
	MaxOpenConnections int
}

type DB struct {
	db *sqlx.DB
}

func Connect(cfg Config) (*DB, error) {
	if cfg.Driver != "postgres" { // We support only postgres right now.
		return nil, fmt.Errorf("unsuported sql driver: %s", cfg.Driver)
	}

	db, err := sql.Open(cfg.Driver, cfg.DSN)
	if err != nil {
		return nil, err
	}
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	if cfg.MaxIdleConnections > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConnections)
	}
	if cfg.MaxOpenConnections > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConnections)
	}
	if err := pingDatabase(db); err != nil {
		return nil, err
	}

	return &DB{db: sqlx.NewDb(db, cfg.Driver)}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func pingDatabase(db *sql.DB) (err error) {
	for i := 0; i < 10; i++ {
		err = db.Ping()
		if err == nil {
			return
		}
		time.Sleep(time.Second)
	}
	return
}
