package payments

import (
	"context"
	"log"

	"github.com/michaljemala/payments-sample/pkg/internal/errors"

	"github.com/michaljemala/payments-sample/pkg/domain"
	"github.com/michaljemala/payments-sample/pkg/internal/service"
	"github.com/michaljemala/payments-sample/pkg/internal/store"
)

type defaultPaymentService struct {
	*service.Generic

	paymentStore paymentStore
	enumStore    enumStore

	logger *log.Logger
}

type (
	paymentStore interface {
		Count(store.Tx, domain.PaymentSearchRequest) (uint, error)
		Find(store.Tx, domain.PaymentSearchRequest) ([]*domain.Payment, error)
		Get(store.Tx, domain.ID) (*domain.Payment, error)
		Insert(store.Tx, *domain.Payment) error
		Delete(store.Tx, domain.ID) error
		Update(store.Tx, *domain.Payment) error
	}
	enumStore interface {
		Exists(tx store.Tx, name domain.EnumName, code string) (bool, error)
	}
)

func newPaymentService(txManager store.TxManager, paymentStore paymentStore, enumStore enumStore, logger *log.Logger) paymentService {
	return &defaultPaymentService{
		Generic:      &service.Generic{TxManager: txManager},
		paymentStore: paymentStore,
		enumStore:    enumStore,
		logger:       logger,
	}
}

func (s *defaultPaymentService) Search(ctx context.Context, searchReq domain.PaymentSearchRequest) (*domain.PaymentSearchResponse, error) {
	searchResp := new(domain.PaymentSearchResponse)
	err := s.WithTransaction(ctx, func(tx store.Tx) error {
		var err error
		searchResp.Data, err = s.paymentStore.Find(tx, searchReq)
		if err != nil {
			return err
		}
		if searchReq.SearchPagination != nil {
			searchResp.Size, err = s.paymentStore.Count(tx, searchReq)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return searchResp, nil
}

func (s *defaultPaymentService) Load(ctx context.Context, id domain.ID) (payment *domain.Payment, err error) {
	err = s.WithTransaction(ctx, func(tx store.Tx) error {
		payment, err = s.paymentStore.Get(tx, id)
		return err
	})
	return payment, err
}

func (s *defaultPaymentService) Create(ctx context.Context, payment *domain.Payment) error {
	return s.WithTransaction(ctx, func(tx store.Tx) error {
		err := s.validatePayment(tx, payment)
		if err != nil {
			return err
		}

		return s.paymentStore.Insert(tx, payment)
	})
}

func (s *defaultPaymentService) Delete(ctx context.Context, id domain.ID) error {
	return s.WithTransaction(ctx, func(tx store.Tx) error {
		return s.paymentStore.Delete(tx, id)
	})
}

func (s *defaultPaymentService) Update(ctx context.Context, payment *domain.Payment) error {
	return s.WithTransaction(ctx, func(tx store.Tx) error {
		err := s.validatePayment(tx, payment)
		if err != nil {
			return err
		}

		return s.paymentStore.Update(tx, payment)
	})
}

func (s *defaultPaymentService) validatePayment(tx store.Tx, payment *domain.Payment) error {
	if payment == nil {
		return errors.Generic(errors.ErrCodeGenericInvalidArgument, "payment must not be nil", "")
	}

	v := &paymentValidator{s.enumStore}

	err := v.enumExists(tx, enumNameCurrency, payment.Amount.Currency, "payment.amount.currency")
	if err != nil {
		return err
	}
	err = v.enumExists(tx, enumNameCountry, payment.Creditor.Address.CountryCode, "payment.creditor.address.country_code")
	if err != nil {
		return err
	}
	err = v.enumExists(tx, enumNameCountry, payment.Debtor.Address.CountryCode, "payment.debtor.address.country_code")
	if err != nil {
		return err
	}
	err = v.enumExists(tx, enumNameScheme, payment.Scheme, "payment.scheme")
	if err != nil {
		return err
	}

	return nil
}

type paymentValidator struct {
	enumStore enumStore
}

func (v *paymentValidator) enumExists(tx store.Tx, name domain.EnumName, code, field string) error {
	ok, err := v.enumStore.Exists(tx, name, code)
	if err != nil {
		return err
	}
	if !ok {
		return errors.Generic(errors.ErrCodeGenericInvalidArgument, "enum not found", field)
	}
	return nil
}
