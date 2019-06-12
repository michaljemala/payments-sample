package store

import "context"

type TxManager interface {
	Begin(ctx context.Context) (Tx, error)
}

type Tx interface {
	Commit() error
	Rollback() error
}

type Transactional interface {
	WithTransaction(context.Context, func(Tx) error) error
}
