package domain

import (
	"github.com/michaljemala/payments-sample/pkg/internal/errors"
	"github.com/michaljemala/payments-sample/pkg/internal/resource"
)

type Payment struct {
	BaseObject

	Amount   Monetary     `json:"amount"`
	Creditor PaymentParty `json:"creditor"`
	Debtor   PaymentParty `json:"debtor"`
	Scheme   string       `json:"scheme"`
}

func (p Payment) Validate() error {
	if p.ID.IsNil() {
		return errors.Generic(
			errors.ErrCodeGenericInvalidArgument,
			"invalid payment",
			"payment id must not be nil",
		)
	}
	if p.Scheme == "" {
		return errors.Generic(
			errors.ErrCodeGenericInvalidArgument,
			"invalid payment",
			"invalid payment scheme",
		)
	}
	if p.Amount.Currency == "" {
		return errors.Generic(
			errors.ErrCodeGenericInvalidArgument,
			"invalid payment",
			"invalid amount currency code",
		)
	}
	if p.Creditor.AccountNumber == "" {
		return errors.Generic(
			errors.ErrCodeGenericInvalidArgument,
			"invalid payment",
			"invalid creditor account number",
		)
	}
	if p.Debtor.AccountNumber == "" {
		return errors.Generic(
			errors.ErrCodeGenericInvalidArgument,
			"invalid payment",
			"invalid debtor account number",
		)
	}
	return nil
}

type PaymentParty struct {
	Name            string          `json:"name"`
	Address         Address         `json:"address"`
	AccountName     string          `json:"account_name"`
	AccountNumber   string          `json:"account_number"`
	AccountProvider AccountProvider `json:"account_provider"`
}

type AccountProvider struct {
	Code string  `json:"code"`
	Name *string `json:"name,omitempty"`
}

type Monetary struct {
	Value    Decimal `json:"value"`
	Currency string  `json:"currency"`
}

type PaymentSearchRequest struct {
	*resource.SearchPagination
	resource.SearchFilter
}

func (r PaymentSearchRequest) IDs() []ID {
	if r.SearchFilter == nil {
		return nil
	}
	ids, ok := r.SearchFilter["id"].([]ID)
	if !ok {
		return nil
	}
	return ids
}

func (r PaymentSearchRequest) CreditorAccountNumbers() []string {
	if r.SearchFilter == nil {
		return nil
	}
	numbers, ok := r.SearchFilter["creditor.account_number"].([]string)
	if !ok {
		return nil
	}
	return numbers
}

func (r PaymentSearchRequest) DebtorAccountNumbers() []string {
	if r.SearchFilter == nil {
		return nil
	}
	numbers, ok := r.SearchFilter["debtor.account_number"].([]string)
	if !ok {
		return nil
	}
	return numbers
}

type PaymentSearchResponse struct {
	Data []*Payment
	Size uint
}
