package sql

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/michaljemala/payments-sample/pkg/internal/store"
)

type TxManager struct {
	db *DB
}

func NewTxManager(db *DB) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) Begin(ctx context.Context) (store.Tx, error) {
	txx, err := m.db.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &Tx{tx: txx}, nil
}

type Tx struct {
	tx *sqlx.Tx
}

func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

func (tx *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return tx.tx.Query(tx.tx.Rebind(query), args...)
}

func (tx *Tx) QueryRow(query string, args ...interface{}) *sql.Row {
	return tx.tx.QueryRow(tx.tx.Rebind(query), args...)
}

func (tx *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.tx.Exec(tx.tx.Rebind(query), args...)
}
