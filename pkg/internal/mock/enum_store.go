package mock

import (
	"github.com/michaljemala/payments-sample/pkg/domain"
	"github.com/michaljemala/payments-sample/pkg/internal/store"
)

type EnumStore struct {
	ExistsFn      func(store.Tx, domain.EnumName, string) (bool, error)
	ExistsInvoked bool
}

func (s *EnumStore) Exists(tx store.Tx, name domain.EnumName, code string) (bool, error) {
	s.ExistsInvoked = true
	return s.ExistsFn(tx, name, code)
}
