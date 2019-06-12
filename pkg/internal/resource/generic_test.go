package resource

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/michaljemala/payments-sample/pkg/internal/errors"
)

func TestGeneric_ExtractSearchFilter(t *testing.T) {
	testCases := []struct {
		name         string
		paramFunc    func(string, []string) (interface{}, error)
		filterParams map[string][]string
		errFunc      func(*testing.T, error)
		searchFilter SearchFilter
	}{
		{
			name: "Supported filter parameter",
			paramFunc: func(key string, values []string) (interface{}, error) {
				if key == "foo" {
					return values, nil
				}
				return nil, errors.Generic(errors.ErrCodeGenericInvalidArgument, "invalid search filter", "")
			},
			filterParams: map[string][]string{
				"filter[foo]": {"bar"},
			},
			searchFilter: SearchFilter{
				"foo": []string{"bar"},
			},
		},
		{
			name: "Unsupported filter parameter",
			paramFunc: func(key string, values []string) (interface{}, error) {
				return nil, fmt.Errorf("some error")
			},
			filterParams: map[string][]string{
				"filter[UNKNOWN]": {"bar"},
			},
			errFunc: func(t *testing.T, err error) {
				switch err := err.(type) {
				case errors.Error:
					if want, have := errors.ErrCodeGenericInvalidArgument, err.Code; want != have {
						t.Fatalf("unexpected error code: want %s, have %s", want, have)
					}
				case nil:
					t.Fatal("expected error, have <nil")
				default:
					t.Fatalf("unexpected error: (%T)%v", err, err)
				}
			},
		},
		{
			name: "Invalid filter parameter",
			paramFunc: func(key string, values []string) (interface{}, error) {
				return nil, fmt.Errorf("some error")
			},
			filterParams: map[string][]string{
				"I_AM_NOT_VALID": {"OOOPS"},
			},
			errFunc: func(t *testing.T, err error) {
				switch err := err.(type) {
				case errors.Error:
					if want, have := errors.ErrCodeGenericInvalidArgument, err.Code; want != have {
						t.Fatalf("unexpected error code: want %s, have %s", want, have)
					}
				case nil:
					t.Fatal("expected error, have <nil")
				default:
					t.Fatalf("unexpected error: (%T)%v", err, err)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := &Generic{ParamFunc: tc.paramFunc}
			searchFilter, err := g.ExtractSearchFilter(tc.filterParams)
			if err != nil {
				if tc.errFunc == nil {
					t.Fatal("unexpected error: errFunc not provided")
				}
				tc.errFunc(t, err)
			}
			if want, have := tc.searchFilter, searchFilter; !cmp.Equal(want, have) {
				t.Fatalf("invalid search filter: %v", cmp.Diff(want, have))
			}
		})
	}
}

func TestGeneric_ExtractPagination(t *testing.T) {
	testCases := []struct {
		name             string
		paginationParams map[string]string
		errFunc          func(*testing.T, error)
		pagination       *SearchPagination
	}{
		{
			name: "Valid pagination parameters",
			paginationParams: map[string]string{
				"number": "1",
				"size":   "100",
			},
			pagination: &SearchPagination{
				page: 1,
				size: 100,
			},
		},
		{
			name: "Unsupported pagination parameters",
			paginationParams: map[string]string{
				"UNSUPPORTED": "1000",
			},
			errFunc: func(t *testing.T, err error) {
				switch err := err.(type) {
				case errors.Error:
					if want, have := errors.ErrCodeGenericInvalidArgument, err.Code; want != have {
						t.Fatalf("unexpected error code: want %s, have %s", want, have)
					}
				case nil:
					t.Fatal("expected error, have <nil")
				default:
					t.Fatalf("unexpected error: (%T)%v", err, err)
				}
			},
		},
		{
			name: "Invalid pagination parameters",
			paginationParams: map[string]string{
				"number": "THIS_SHOULD_BE_A_UINT",
			},
			errFunc: func(t *testing.T, err error) {
				switch err := err.(type) {
				case errors.Error:
					if want, have := errors.ErrCodeGenericInvalidArgument, err.Code; want != have {
						t.Fatalf("unexpected error code: want %s, have %s", want, have)
					}
				case nil:
					t.Fatal("expected error, have <nil")
				default:
					t.Fatalf("unexpected error: (%T)%v", err, err)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			g := &Generic{}
			pagination, err := g.ExtractPagination(tc.paginationParams)
			if err != nil {
				if tc.errFunc == nil {
					t.Fatal("unexpected error: errFunc not provided")
				}
				tc.errFunc(t, err)
			}
			opts := cmp.Options{
				cmp.AllowUnexported(SearchPagination{}),
			}
			if want, have := tc.pagination, pagination; !cmp.Equal(want, have, opts...) {
				t.Fatalf("invalid pagination: %v", cmp.Diff(want, have, opts...))
			}
		})
	}

}
