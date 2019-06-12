package domain

import (
	"testing"

	"github.com/michaljemala/payments-sample/pkg/internal/errors"
)

func TestPayment_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		in      Payment
		errFunc func(t *testing.T, err error)
	}{
		{
			name: "No ID",
			in: Payment{
				Amount: Monetary{
					Value:    MustDecimalFrom("1000.0"),
					Currency: "EUR",
				},
				Creditor: PaymentParty{
					AccountNumber: "SK3302000000000000012351",
				},
				Debtor: PaymentParty{
					AccountNumber: "SK0809000000000123123123",
				},
				Scheme: "SWIFT",
			},
			errFunc: assertInvalidArgumentError,
		},
		{
			name: "No amount currency",
			in: Payment{
				BaseObject: BaseObject{ID: MustIDFrom("276c8bbf-79ca-4ac2-b319-0f1c51463540")},
				Amount: Monetary{
					Value: MustDecimalFrom("1000.0"),
				},
				Creditor: PaymentParty{
					AccountNumber: "SK3302000000000000012351",
				},
				Debtor: PaymentParty{
					AccountNumber: "SK0809000000000123123123",
				},
				Scheme: "SWIFT",
			},
			errFunc: assertInvalidArgumentError,
		},
		{
			name: "No scheme",
			in: Payment{
				BaseObject: BaseObject{ID: MustIDFrom("276c8bbf-79ca-4ac2-b319-0f1c51463540")},
				Amount: Monetary{
					Value:    MustDecimalFrom("1000.0"),
					Currency: "EUR",
				},
				Creditor: PaymentParty{
					AccountNumber: "SK3302000000000000012351",
				},
				Debtor: PaymentParty{
					AccountNumber: "SK0809000000000123123123",
				},
			},
			errFunc: assertInvalidArgumentError,
		},
		{
			name: "No creditor account number",
			in: Payment{
				BaseObject: BaseObject{ID: MustIDFrom("276c8bbf-79ca-4ac2-b319-0f1c51463540")},
				Amount: Monetary{
					Value:    MustDecimalFrom("1000.0"),
					Currency: "EUR",
				},
				Debtor: PaymentParty{
					AccountNumber: "SK0809000000000123123123",
				},
				Scheme: "SWIFT",
			},
			errFunc: assertInvalidArgumentError,
		},
		{
			name: "No debtor account number",
			in: Payment{
				BaseObject: BaseObject{ID: MustIDFrom("276c8bbf-79ca-4ac2-b319-0f1c51463540")},
				Amount: Monetary{
					Value:    MustDecimalFrom("1000.0"),
					Currency: "EUR",
				},
				Creditor: PaymentParty{
					AccountNumber: "SK3302000000000000012351",
				},
				Scheme: "SWIFT",
			},
			errFunc: assertInvalidArgumentError,
		},
		{
			name: "Valid payment",
			in: Payment{
				BaseObject: BaseObject{ID: MustIDFrom("276c8bbf-79ca-4ac2-b319-0f1c51463540")},
				Amount: Monetary{
					Value:    MustDecimalFrom("1000.0"),
					Currency: "EUR",
				},
				Creditor: PaymentParty{
					AccountNumber: "SK3302000000000000012351",
				},
				Debtor: PaymentParty{
					AccountNumber: "SK0809000000000123123123",
				},
				Scheme: "SWIFT",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.in.Validate()
			if err != nil {
				if tc.errFunc == nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

func assertInvalidArgumentError(t *testing.T, err error) {
	switch err := err.(type) {
	case errors.Error:
		if want, have := errors.ErrCodeGenericInvalidArgument, err.Code; want != have {
			t.Fatalf("invalid error code: want %s, have %s", want, have)
		}
	default:
		t.Fatalf("invalid error type: want %T, have %T", errors.Error{}, err)
	}
}
