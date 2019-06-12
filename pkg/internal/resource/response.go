package resource

import (
	"fmt"
	"net/http"

	"github.com/manyminds/api2go"

	"github.com/michaljemala/payments-sample/pkg/internal/errors"
)

func WrapError(err error) api2go.HTTPError {
	translated, status := translateError(err)
	httpErr := api2go.NewHTTPError(err, "", status)
	httpErr.Errors = append(httpErr.Errors, translated...)
	return httpErr
}

func translateError(err error) ([]api2go.Error, int) {
	switch err := err.(type) {
	case errors.Error:
		switch err.Category {
		case errors.ErrCategoryGeneric:
			switch err.Code {
			case errors.ErrCodeGenericAlreadyExists:
				return []api2go.Error{{
					Status: fmt.Sprintf("%d", http.StatusConflict),
					Code:   err.Code.String(),
					Title:  err.Message,
					Detail: err.Detail,
				}}, http.StatusConflict
			case errors.ErrCodeGenericInvalidArgument:
				return []api2go.Error{{
					Status: fmt.Sprintf("%d", http.StatusBadRequest),
					Code:   err.Code.String(),
					Title:  err.Message,
					Detail: err.Detail,
				}}, http.StatusBadRequest
			case errors.ErrCodeGenericNotFound:
				return []api2go.Error{{
					Status: fmt.Sprintf("%d", http.StatusNotFound),
					Code:   err.Code.String(),
					Title:  err.Message,
					Detail: err.Detail,
				}}, http.StatusNotFound
			default:
				return []api2go.Error{{
					Status: fmt.Sprintf("%d", http.StatusInternalServerError),
					Code:   err.Code.String(),
					Title:  err.Message,
					Detail: err.Detail,
				}}, http.StatusInternalServerError
			}
		}
	}

	return []api2go.Error{{
		Status: fmt.Sprintf("%d", http.StatusInternalServerError),
		Code:   errors.ErrCodeGenericInternal.String(),
	}}, http.StatusInternalServerError
}

func WrapObject(v interface{}, status int) api2go.Responder {
	return &api2go.Response{Res: v, Code: status}
}

func WrapArray(v interface{}, status int) api2go.Responder {
	return &api2go.Response{Res: v, Code: status}
}
