package sql

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/lib/pq"
	"github.com/michaljemala/pqerror"

	"github.com/michaljemala/payments-sample/pkg/internal/errors"
)

func TestWrapSelectError(t *testing.T) {
	testCases := []struct {
		name    string
		in      error
		errFunc func(*testing.T, error)
	}{
		{
			name: "No error",
			errFunc: func(t *testing.T, err error) {
				if err != nil {
					t.Fatalf("unexpected error: want <nil>, have %v", err)
				}
			},
		},
		{
			name: "No rows",
			in:   sql.ErrNoRows,
			errFunc: func(t *testing.T, err error) {
				switch err := err.(type) {
				case errors.Error:
					if want, have := errors.ErrCodeGenericNotFound, err.Code; want != have {
						t.Fatalf("unexpected error code: want %s, have %s", want, have)
					}
				default:
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},
		{
			name: "Unknown error",
			in:   fmt.Errorf("I_AM_UNKNOWN"),
			errFunc: func(t *testing.T, err error) {
				switch err := err.(type) {
				case errors.Error:
					if want, have := errors.ErrCodeDataAccessSelectFailed, err.Code; want != have {
						t.Fatalf("unexpected error code: want %s, have %s", want, have)
					}
				default:
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := WrapSelectError(tc.in, "")
			if tc.errFunc != nil {
				tc.errFunc(t, err)
			}
		})
	}
}

func TestWrapInsertError(t *testing.T) {
	testCases := []struct {
		name    string
		in      error
		errFunc func(*testing.T, error)
	}{
		{
			name: "No error",
			errFunc: func(t *testing.T, err error) {
				if err != nil {
					t.Fatalf("unexpected error: want <nil>, have %v", err)
				}
			},
		},
		{
			name: "Conflict",
			in:   &pq.Error{Code: pqerror.UniqueViolation},
			errFunc: func(t *testing.T, err error) {
				switch err := err.(type) {
				case errors.Error:
					if want, have := errors.ErrCodeGenericAlreadyExists, err.Code; want != have {
						t.Fatalf("unexpected error code: want %s, have %s", want, have)
					}
				default:
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},
		{
			name: "Unknown error",
			in:   fmt.Errorf("I_AM_UNKNOWN"),
			errFunc: func(t *testing.T, err error) {
				switch err := err.(type) {
				case errors.Error:
					if want, have := errors.ErrCodeDataAccessInsertFailed, err.Code; want != have {
						t.Fatalf("unexpected error code: want %s, have %s", want, have)
					}
				default:
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := WrapInsertError(tc.in, "")
			if tc.errFunc != nil {
				tc.errFunc(t, err)
			}
		})
	}
}

func TestWrapDeleteError(t *testing.T) {
	testCases := []struct {
		name    string
		in      error
		errFunc func(*testing.T, error)
	}{
		{
			name: "No error",
			errFunc: func(t *testing.T, err error) {
				if err != nil {
					t.Fatalf("unexpected error: want <nil>, have %v", err)
				}
			},
		},
		{
			name: "No rows",
			in:   sql.ErrNoRows,
			errFunc: func(t *testing.T, err error) {
				switch err := err.(type) {
				case errors.Error:
					if want, have := errors.ErrCodeGenericNotFound, err.Code; want != have {
						t.Fatalf("unexpected error code: want %s, have %s", want, have)
					}
				default:
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},
		{
			name: "Unknown error",
			in:   fmt.Errorf("I_AM_UNKNOWN"),
			errFunc: func(t *testing.T, err error) {
				switch err := err.(type) {
				case errors.Error:
					if want, have := errors.ErrCodeDataAccessDeleteFailed, err.Code; want != have {
						t.Fatalf("unexpected error code: want %s, have %s", want, have)
					}
				default:
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := WrapDeleteError(tc.in, "")
			if tc.errFunc != nil {
				tc.errFunc(t, err)
			}
		})
	}
}

func TestWrapUpdateError(t *testing.T) {
	testCases := []struct {
		name    string
		in      error
		errFunc func(*testing.T, error)
	}{
		{
			name: "No error",
			errFunc: func(t *testing.T, err error) {
				if err != nil {
					t.Fatalf("unexpected error: want <nil>, have %v", err)
				}
			},
		},
		{
			name: "No rows",
			in:   sql.ErrNoRows,
			errFunc: func(t *testing.T, err error) {
				switch err := err.(type) {
				case errors.Error:
					if want, have := errors.ErrCodeGenericNotFound, err.Code; want != have {
						t.Fatalf("unexpected error code: want %s, have %s", want, have)
					}
				default:
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},
		{
			name: "Unknown error",
			in:   fmt.Errorf("I_AM_UNKNOWN"),
			errFunc: func(t *testing.T, err error) {
				switch err := err.(type) {
				case errors.Error:
					if want, have := errors.ErrCodeDataAccessUpdateFailed, err.Code; want != have {
						t.Fatalf("unexpected error code: want %s, have %s", want, have)
					}
				default:
					t.Fatalf("unexpected error: %v", err)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := WrapUpdateError(tc.in, "")
			if tc.errFunc != nil {
				tc.errFunc(t, err)
			}
		})
	}
}
