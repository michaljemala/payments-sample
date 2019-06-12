package payments

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/manyminds/api2go/jsonapi"

	"github.com/michaljemala/payments-sample/pkg/domain"
	"github.com/michaljemala/payments-sample/pkg/internal/errors"
	"github.com/michaljemala/payments-sample/pkg/internal/mock"
	"github.com/michaljemala/payments-sample/pkg/internal/service"
	"github.com/michaljemala/payments-sample/pkg/internal/store"
)

func TestMain(m *testing.M) {
	log.SetOutput(ioutil.Discard) //api2go uses global logger, hide those messages
	os.Exit(m.Run())
}

func TestPayment_Create(t *testing.T) {
	testCases := []struct {
		name         string
		paymentStore paymentStore
		enumStore    enumStore
		in           domain.Payment
		reqFunc      func(*testing.T, *http.Request)
		statusCode   int
		out          domain.Payment
		respFunc     func(*testing.T, *http.Response)
	}{
		{
			name: "Valid payment",
			paymentStore: &mock.PaymentStore{
				InsertFn: func(store.Tx, *domain.Payment) error { return nil },
			},
			in: domain.Payment{
				BaseObject: domain.BaseObject{ID: domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d")},
				Scheme:     "SWIFT",
				Amount: domain.Monetary{
					Value:    domain.MustDecimalFrom("100.00"),
					Currency: "EUR",
				},
				Debtor:   domain.PaymentParty{AccountNumber: "0123456789"},
				Creditor: domain.PaymentParty{AccountNumber: "9876543210"},
			},
			statusCode: http.StatusCreated,
			out: domain.Payment{
				BaseObject: domain.BaseObject{ID: domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d")},
				Scheme:     "SWIFT",
				Amount: domain.Monetary{
					Value:    domain.MustDecimalFrom("100.00"),
					Currency: "EUR",
				},
				Debtor:   domain.PaymentParty{AccountNumber: "0123456789"},
				Creditor: domain.PaymentParty{AccountNumber: "9876543210"},
			},
		},
		{
			name: "Conflicting payment",
			paymentStore: &mock.PaymentStore{
				InsertFn: func(store.Tx, *domain.Payment) error {
					return errors.Generic(errors.ErrCodeGenericAlreadyExists, "payment already exists", "")
				},
			},
			in: domain.Payment{
				BaseObject: domain.BaseObject{ID: domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d")},
				Scheme:     "SWIFT",
				Amount: domain.Monetary{
					Value:    domain.MustDecimalFrom("100.00"),
					Currency: "EUR",
				},
				Debtor:   domain.PaymentParty{AccountNumber: "0123456789"},
				Creditor: domain.PaymentParty{AccountNumber: "9876543210"},
			},
			respFunc: func(t *testing.T, resp *http.Response) {
				if want, have := http.StatusConflict, resp.StatusCode; want != have {
					t.Fatalf("unexpected response status: want %d, have %d", want, have)
				}
			},
		},
		{
			name: "Invalid payment",
			in:   domain.Payment{},
			respFunc: func(t *testing.T, resp *http.Response) {
				if want, have := http.StatusBadRequest, resp.StatusCode; want != have {
					t.Fatalf("unexpected response status: want %d, have %d", want, have)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, close := testPaymentHandler(t, tc.paymentStore, tc.enumStore)
			defer close()

			body, err := jsonapi.Marshal(tc.in)
			if err != nil {
				t.Fatalf("unable to marshal json api payload: %v", err)
			}

			req, err := http.NewRequest("POST", "/payments", bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("unable to create request: %v", err)
			}

			if tc.reqFunc != nil {
				tc.reqFunc(t, req)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			resp := rec.Result()

			if tc.respFunc != nil {
				tc.respFunc(t, resp)
				return
			}

			if want, have := tc.statusCode, resp.StatusCode; want != have {
				t.Fatalf("invalid response status: want %v, have %v", want, have)
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("unable to read response body: %v", err)
			}

			var out domain.Payment
			err = jsonapi.Unmarshal(data, &out)
			if err != nil {
				t.Fatalf("unable to unmarshal json api payload: %v", err)
			}

			opts := []cmp.Option{
				cmp.Transformer("Decimal", func(in domain.Decimal) string {
					return in.String()
				}),
			}

			if want, have := tc.out, out; !cmp.Equal(want, have, opts...) {
				t.Fatalf("invalid payment: %v", cmp.Diff(want, have, opts...))
			}
		})
	}
}

func TestPayment_FindOne(t *testing.T) {
	testCases := []struct {
		name         string
		paymentStore paymentStore
		enumStore    enumStore
		in           domain.ID
		reqFunc      func(*testing.T, *http.Request)
		statusCode   int
		out          domain.Payment
		respFunc     func(*testing.T, *http.Response)
	}{
		{
			name: "Existing payment",
			paymentStore: &mock.PaymentStore{
				GetFn: func(tx store.Tx, id domain.ID) (*domain.Payment, error) {
					return &domain.Payment{
						BaseObject: domain.BaseObject{ID: domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d")},
						Scheme:     "SWIFT",
						Amount: domain.Monetary{
							Value:    domain.MustDecimalFrom("100.00"),
							Currency: "EUR",
						},
						Debtor:   domain.PaymentParty{AccountNumber: "0123456789"},
						Creditor: domain.PaymentParty{AccountNumber: "9876543210"},
					}, nil
				},
			},
			in:         domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d"),
			statusCode: http.StatusOK,
			out: domain.Payment{
				BaseObject: domain.BaseObject{ID: domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d")},
				Scheme:     "SWIFT",
				Amount: domain.Monetary{
					Value:    domain.MustDecimalFrom("100.00"),
					Currency: "EUR",
				},
				Debtor:   domain.PaymentParty{AccountNumber: "0123456789"},
				Creditor: domain.PaymentParty{AccountNumber: "9876543210"},
			},
		},
		{
			name: "Missing payment",
			paymentStore: &mock.PaymentStore{
				GetFn: func(tx store.Tx, id domain.ID) (*domain.Payment, error) {
					if id != domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d") {
						t.Fatal("unexpected ID")
					}
					return nil, errors.Generic(errors.ErrCodeGenericNotFound, "payment not found ", "")
				},
			},
			in:         domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d"),
			statusCode: http.StatusOK,
			respFunc: func(t *testing.T, resp *http.Response) {
				if want, have := http.StatusNotFound, resp.StatusCode; want != have {
					t.Fatalf("unexpected response status: want %d, have %d", want, have)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, close := testPaymentHandler(t, tc.paymentStore, tc.enumStore)
			defer close()

			url := fmt.Sprintf("/payments/%s", tc.in)

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatalf("unable to create request: %v", err)
			}

			if tc.reqFunc != nil {
				tc.reqFunc(t, req)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			resp := rec.Result()

			if tc.respFunc != nil {
				tc.respFunc(t, resp)
				return
			}

			if want, have := tc.statusCode, resp.StatusCode; want != have {
				t.Fatalf("invalid response status: want %v, have %v", want, have)
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("unable to read response body: %v", err)
			}

			var out domain.Payment
			err = jsonapi.Unmarshal(data, &out)
			if err != nil {
				t.Fatalf("unable to unmarshal json api payload: %v", err)
			}

			opts := []cmp.Option{
				cmp.Transformer("Decimal", func(in domain.Decimal) string {
					return in.String()
				}),
			}

			if want, have := tc.out, out; !cmp.Equal(want, have, opts...) {
				t.Fatalf("invalid payment: %v", cmp.Diff(want, have, opts...))
			}
		})
	}
}

func TestPayment_FindAll(t *testing.T) {
	testCases := []struct {
		name         string
		paymentStore paymentStore
		enumStore    enumStore
		in           string
		reqFunc      func(*testing.T, *http.Request)
		statusCode   int
		out          []domain.Payment
		respFunc     func(*testing.T, *http.Response)
	}{
		{
			name: "No search filter",
			paymentStore: &mock.PaymentStore{
				FindFn: func(store.Tx, domain.PaymentSearchRequest) ([]*domain.Payment, error) {
					return []*domain.Payment{
						{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d")}},
						{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("b12fc840-2511-452a-8cdf-407c09eba168")}},
						{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("3ebbc3b4-d1ec-4316-a191-da97bcaec65d")}},
						{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("1caef84f-50a7-446c-bb16-c0b197bba112")}},
						{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("1166b097-9010-42cd-84fe-f843a37d805a")}},
					}, nil
				},
			},
			statusCode: http.StatusOK,
			out: []domain.Payment{
				{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d")}},
				{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("b12fc840-2511-452a-8cdf-407c09eba168")}},
				{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("3ebbc3b4-d1ec-4316-a191-da97bcaec65d")}},
				{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("1caef84f-50a7-446c-bb16-c0b197bba112")}},
				{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("1166b097-9010-42cd-84fe-f843a37d805a")}},
			},
		},
		{
			name: "Valid search filter",
			paymentStore: &mock.PaymentStore{
				FindFn: func(tx store.Tx, r domain.PaymentSearchRequest) ([]*domain.Payment, error) {
					if len(r.IDs()) != 2 {
						t.Fatal("invalid search filter")
					}
					if r.IDs()[0] != domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d") {
						t.Fatal("unexpected search filter")
					}
					if r.IDs()[1] != domain.MustIDFrom("b12fc840-2511-452a-8cdf-407c09eba168") {
						t.Fatal("unexpected search filter")
					}
					return []*domain.Payment{
						{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d")}},
						{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("b12fc840-2511-452a-8cdf-407c09eba168")}},
					}, nil
				},
			},
			in:         "filter[id]=33b5c07b-c6bd-4a59-b02b-554256eaba5d,b12fc840-2511-452a-8cdf-407c09eba168",
			statusCode: http.StatusOK,
			out: []domain.Payment{
				{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d")}},
				{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("b12fc840-2511-452a-8cdf-407c09eba168")}},
			},
		},
		{
			name: "Invalid search filter",
			in:   "filter[UNKNOWN]=SOME_VALUE",
			respFunc: func(t *testing.T, resp *http.Response) {
				if want, have := http.StatusBadRequest, resp.StatusCode; want != have {
					t.Fatalf("unexpected response status: want %d, have %d", want, have)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, close := testPaymentHandler(t, tc.paymentStore, tc.enumStore)
			defer close()

			url := fmt.Sprintf("/payments?%s", tc.in)

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatalf("unable to create request: %v", err)
			}

			if tc.reqFunc != nil {
				tc.reqFunc(t, req)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			resp := rec.Result()

			if tc.respFunc != nil {
				tc.respFunc(t, resp)
				return
			}

			if want, have := tc.statusCode, resp.StatusCode; want != have {
				t.Fatalf("invalid response status: want %v, have %v", want, have)
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("unable to read response body: %v", err)
			}

			var out []domain.Payment
			err = jsonapi.Unmarshal(data, &out)
			if err != nil {
				t.Fatalf("unable to unmarshal json api payload: %v", err)
			}

			opts := []cmp.Option{
				cmp.Transformer("Decimal", func(in domain.Decimal) string {
					return in.String()
				}),
			}

			if want, have := tc.out, out; !cmp.Equal(want, have, opts...) {
				t.Fatalf("invalid payment: %v", cmp.Diff(want, have, opts...))
			}
		})
	}
}

func TestPayment_PaginatedFindAll(t *testing.T) {
	testCases := []struct {
		name         string
		paymentStore paymentStore
		enumStore    enumStore
		in           string
		reqFunc      func(*testing.T, *http.Request)
		statusCode   int
		out          []domain.Payment
		respFunc     func(*testing.T, *http.Response)
	}{
		{
			name: "No pagination",
			paymentStore: &mock.PaymentStore{
				CountFn: func(store.Tx, domain.PaymentSearchRequest) (uint, error) {
					return 5, nil
				},
				FindFn: func(store.Tx, domain.PaymentSearchRequest) ([]*domain.Payment, error) {
					return []*domain.Payment{
						{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d")}},
						{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("b12fc840-2511-452a-8cdf-407c09eba168")}},
						{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("3ebbc3b4-d1ec-4316-a191-da97bcaec65d")}},
						{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("1caef84f-50a7-446c-bb16-c0b197bba112")}},
						{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("1166b097-9010-42cd-84fe-f843a37d805a")}},
					}, nil
				},
			},
			statusCode: http.StatusOK,
			out: []domain.Payment{
				{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d")}},
				{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("b12fc840-2511-452a-8cdf-407c09eba168")}},
				{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("3ebbc3b4-d1ec-4316-a191-da97bcaec65d")}},
				{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("1caef84f-50a7-446c-bb16-c0b197bba112")}},
				{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("1166b097-9010-42cd-84fe-f843a37d805a")}},
			},
		},
		{
			name: "Valid pagination",
			paymentStore: &mock.PaymentStore{
				CountFn: func(store.Tx, domain.PaymentSearchRequest) (uint, error) {
					return 5, nil
				},
				FindFn: func(tx store.Tx, r domain.PaymentSearchRequest) ([]*domain.Payment, error) {
					if r.SearchPagination == nil {
						t.Fatal("invalid pagination")
					}
					if r.SearchPagination.Page() != 1 {
						t.Fatal("unexpected page number")
					}
					if r.SearchPagination.Size() != 2 {
						t.Fatal("unexpected page size")
					}
					return []*domain.Payment{
						{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d")}},
						{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("b12fc840-2511-452a-8cdf-407c09eba168")}},
					}, nil
				},
			},
			in:         "page[number]=1&page[size]=2",
			statusCode: http.StatusOK,
			out: []domain.Payment{
				{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d")}},
				{BaseObject: domain.BaseObject{ID: domain.MustIDFrom("b12fc840-2511-452a-8cdf-407c09eba168")}},
			},
		},
		{
			name: "Invalid pagination",
			in:   "page[UNKNOWN]=SOME_VALUE",
			respFunc: func(t *testing.T, resp *http.Response) {
				if want, have := http.StatusBadRequest, resp.StatusCode; want != have {
					t.Fatalf("unexpected response status: want %d, have %d", want, have)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, close := testPaymentHandler(t, tc.paymentStore, tc.enumStore)
			defer close()

			url := fmt.Sprintf("/payments?%s", tc.in)

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatalf("unable to create request: %v", err)
			}

			if tc.reqFunc != nil {
				tc.reqFunc(t, req)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			resp := rec.Result()

			if tc.respFunc != nil {
				tc.respFunc(t, resp)
				return
			}

			if want, have := tc.statusCode, resp.StatusCode; want != have {
				t.Fatalf("invalid response status: want %v, have %v", want, have)
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("unable to read response body: %v", err)
			}

			var out []domain.Payment
			err = jsonapi.Unmarshal(data, &out)
			if err != nil {
				t.Fatalf("unable to unmarshal json api payload: %v", err)
			}

			opts := []cmp.Option{
				cmp.Transformer("Decimal", func(in domain.Decimal) string {
					return in.String()
				}),
			}

			if want, have := tc.out, out; !cmp.Equal(want, have, opts...) {
				t.Fatalf("invalid payment: %v", cmp.Diff(want, have, opts...))
			}
		})
	}
}

func TestPayment_Update(t *testing.T) {
	testCases := []struct {
		name         string
		paymentStore paymentStore
		enumStore    enumStore
		in           domain.Payment
		reqFunc      func(*testing.T, *http.Request)
		statusCode   int
		out          domain.Payment
		respFunc     func(*testing.T, *http.Response)
	}{
		{
			name: "Valid payment",
			paymentStore: &mock.PaymentStore{
				GetFn: func(tx store.Tx, id domain.ID) (*domain.Payment, error) {
					return &domain.Payment{
						BaseObject: domain.BaseObject{ID: domain.MustIDFrom("b12fc840-2511-452a-8cdf-407c09eba168")},
						Amount: domain.Monetary{
							Currency: "EUR",
							Value:    domain.MustDecimalFrom("0"),
						},
					}, nil
				},
				UpdateFn: func(tx store.Tx, p *domain.Payment) error {
					if p == nil {
						t.Fatal("unexpected payment")
					}
					p.Amount.Value = domain.MustDecimalFrom("100.00")
					return nil
				},
			},
			in: domain.Payment{
				BaseObject: domain.BaseObject{ID: domain.MustIDFrom("b12fc840-2511-452a-8cdf-407c09eba168")},
				Amount: domain.Monetary{
					Value:    domain.MustDecimalFrom("100.00"),
					Currency: "EUR",
				},
			},
			statusCode: http.StatusOK,
			out: domain.Payment{
				BaseObject: domain.BaseObject{ID: domain.MustIDFrom("b12fc840-2511-452a-8cdf-407c09eba168")},
				Amount: domain.Monetary{
					Value:    domain.MustDecimalFrom("100.00"),
					Currency: "EUR",
				},
			},
		},
		{
			name: "Missing payment",
			paymentStore: &mock.PaymentStore{
				GetFn: func(tx store.Tx, id domain.ID) (*domain.Payment, error) {
					return nil, errors.Generic(errors.ErrCodeGenericNotFound, "payment not found", "")
				},
			},
			in: domain.Payment{},
			respFunc: func(t *testing.T, resp *http.Response) {
				if want, have := http.StatusNotFound, resp.StatusCode; want != have {
					t.Fatalf("unexpected response status: want %d, have %d", want, have)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, close := testPaymentHandler(t, tc.paymentStore, tc.enumStore)
			defer close()

			body, err := jsonapi.Marshal(tc.in)
			if err != nil {
				t.Fatalf("unable to marshal json api payload: %v", err)
			}

			url := fmt.Sprintf("/payments/%s", tc.in.ID.String())

			req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
			if err != nil {
				t.Fatalf("unable to create request: %v", err)
			}

			if tc.reqFunc != nil {
				tc.reqFunc(t, req)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			resp := rec.Result()

			if tc.respFunc != nil {
				tc.respFunc(t, resp)
				return
			}

			if want, have := tc.statusCode, resp.StatusCode; want != have {
				t.Fatalf("invalid response status: want %v, have %v", want, have)
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("unable to read response body: %v", err)
			}

			var out domain.Payment
			err = jsonapi.Unmarshal(data, &out)
			if err != nil {
				t.Fatalf("unable to unmarshal json api payload: %v", err)
			}

			opts := []cmp.Option{
				cmp.Transformer("Decimal", func(in domain.Decimal) string {
					return in.String()
				}),
			}

			if want, have := tc.out, out; !cmp.Equal(want, have, opts...) {
				t.Fatalf("invalid payment: %v", cmp.Diff(want, have, opts...))
			}
		})
	}
}

func TestPayment_Delete(t *testing.T) {
	testCases := []struct {
		name         string
		paymentStore paymentStore
		enumStore    enumStore
		in           domain.ID
		reqFunc      func(*testing.T, *http.Request)
		statusCode   int
		respFunc     func(*testing.T, *http.Response)
	}{
		{
			name: "Existing payment",
			paymentStore: &mock.PaymentStore{
				DeleteFn: func(tx store.Tx, id domain.ID) error {
					return nil
				},
			},
			in:         domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d"),
			statusCode: http.StatusNoContent,
		},
		{
			name: "Missing payment (idempotent)",
			paymentStore: &mock.PaymentStore{
				DeleteFn: func(tx store.Tx, id domain.ID) error {
					return nil
				},
			},
			in:         domain.MustIDFrom("33b5c07b-c6bd-4a59-b02b-554256eaba5d"),
			statusCode: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler, close := testPaymentHandler(t, tc.paymentStore, tc.enumStore)
			defer close()

			url := fmt.Sprintf("/payments/%s", tc.in)

			req, err := http.NewRequest("DELETE", url, nil)
			if err != nil {
				t.Fatalf("unable to create request: %v", err)
			}

			if tc.reqFunc != nil {
				tc.reqFunc(t, req)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			resp := rec.Result()

			if tc.respFunc != nil {
				tc.respFunc(t, resp)
				return
			}

			if want, have := tc.statusCode, resp.StatusCode; want != have {
				t.Fatalf("invalid response status: want %v, have %v", want, have)
			}

			data, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("unable to read response body: %v", err)
			}
			if len(data) > 0 {
				t.Fatalf("invalid response body: %v", data)
			}
		})
	}
}

func testPaymentHandler(t *testing.T, paymentStore paymentStore, enumStore enumStore) (*API, func()) {
	t.Helper()

	if enumStore == nil {
		enumStore = &mock.EnumStore{
			ExistsFn: func(store.Tx, domain.EnumName, string) (bool, error) {
				return true, nil
			},
		}
	}

	api := newAPI(Config{}, &defaultPaymentService{
		Generic:      &service.Generic{TxManager: &mock.TxManager{}},
		paymentStore: paymentStore,
		enumStore:    enumStore,
	})
	return api, func() {
		err := api.Close()
		if err != nil {
			t.Fatalf("unable to tear down payment handler: %v", err)
		}
	}
}
