package mock

import (
	"context"

	"github.com/michaljemala/payments-sample/pkg/internal/store"
)

type TxManager struct{}

func (m *TxManager) Begin(ctx context.Context) (store.Tx, error) {
	return &Tx{}, nil
}
