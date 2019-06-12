package service

import (
	"context"
	"fmt"

	"github.com/michaljemala/payments-sample/pkg/internal/errors"
	"github.com/michaljemala/payments-sample/pkg/internal/store"
)

type Generic struct {
	TxManager store.TxManager
}

func (s *Generic) WithTransaction(ctx context.Context, fn func(tx store.Tx) error) (err error) {
	if s.TxManager == nil {
		return errors.Generic(
			errors.ErrCodeGenericInternal,
			"Internal Server Error",
			"transaction manager not provided",
		)
	}

	defer func() {
		if r := recover(); r != nil {
			if err == nil {
				switch r := r.(type) {
				case error:
					err = r
				default:
					err = errors.Generic(
						errors.ErrCodeGenericInternal,
						"Internal Server Error",
						fmt.Sprintf("%v", r),
					)
				}
			}
		}
	}()

	var tx store.Tx
	tx, err = s.TxManager.Begin(ctx)
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		panic(tx.Rollback())
	}

	panic(tx.Commit())
}
