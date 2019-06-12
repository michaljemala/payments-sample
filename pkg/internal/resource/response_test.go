package resource

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/michaljemala/payments-sample/pkg/internal/errors"
)

func TestTranslateError(t *testing.T) {
	testCases := []struct {
		name   string
		in     error
		status int
	}{
		{
			name:   "Unexpected error type",
			in:     fmt.Errorf("I_AM_UNEXPECTED"),
			status: http.StatusInternalServerError,
		},
		{
			name:   "Conflict",
			in:     errors.Generic(errors.ErrCodeGenericAlreadyExists, "already exists", ""),
			status: http.StatusConflict,
		},
		{
			name:   "BadRequest",
			in:     errors.Generic(errors.ErrCodeGenericInvalidArgument, "invalid argument", ""),
			status: http.StatusBadRequest,
		},
		{
			name:   "NotFound",
			in:     errors.Generic(errors.ErrCodeGenericNotFound, "not found", ""),
			status: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, status := translateError(tc.in)

			if want, have := tc.status, status; want != have {
				t.Errorf("unexpected message: want %d, have %d", want, have)
			}
		})
	}
}
