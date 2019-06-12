package mock

import (
	"github.com/michaljemala/payments-sample/pkg/domain"
	"github.com/michaljemala/payments-sample/pkg/internal/store"
)

type PaymentStore struct {
	CountFn      func(store.Tx, domain.PaymentSearchRequest) (uint, error)
	CountInvoked bool

	FindFn      func(store.Tx, domain.PaymentSearchRequest) ([]*domain.Payment, error)
	FindInvoked bool

	GetFn      func(store.Tx, domain.ID) (*domain.Payment, error)
	GetInvoked bool

	InsertFn      func(store.Tx, *domain.Payment) error
	InsertInvoked bool

	DeleteFn      func(store.Tx, domain.ID) error
	DeleteInvoked bool

	UpdateFn      func(store.Tx, *domain.Payment) error
	UpdateInvoked bool
}

func (s *PaymentStore) Count(tx store.Tx, r domain.PaymentSearchRequest) (uint, error) {
	s.CountInvoked = true
	return s.CountFn(tx, r)
}

func (s *PaymentStore) Find(tx store.Tx, r domain.PaymentSearchRequest) ([]*domain.Payment, error) {
	s.FindInvoked = true
	return s.FindFn(tx, r)
}

func (s *PaymentStore) Get(tx store.Tx, id domain.ID) (*domain.Payment, error) {
	s.GetInvoked = true
	return s.GetFn(tx, id)
}

func (s *PaymentStore) Insert(tx store.Tx, p *domain.Payment) error {
	s.InsertInvoked = true
	return s.InsertFn(tx, p)
}

func (s *PaymentStore) Delete(tx store.Tx, id domain.ID) error {
	s.DeleteInvoked = true
	return s.DeleteFn(tx, id)
}

func (s *PaymentStore) Update(tx store.Tx, p *domain.Payment) error {
	s.UpdateInvoked = true
	return s.UpdateFn(tx, p)
}
