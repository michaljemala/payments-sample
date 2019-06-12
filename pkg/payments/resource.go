package payments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/manyminds/api2go"

	"github.com/michaljemala/payments-sample/pkg/domain"
	"github.com/michaljemala/payments-sample/pkg/internal/errors"
	"github.com/michaljemala/payments-sample/pkg/internal/resource"
)

type paymentService interface {
	Search(context.Context, domain.PaymentSearchRequest) (*domain.PaymentSearchResponse, error)
	Load(context.Context, domain.ID) (*domain.Payment, error)
	Create(context.Context, *domain.Payment) error
	Delete(context.Context, domain.ID) error
	Update(context.Context, *domain.Payment) error
}

type Resource struct {
	*resource.Generic
	service paymentService
}

func newResource(service paymentService) Resource {
	return Resource{
		Generic: &resource.Generic{
			ParamFunc: paymentParamFunc,
		},
		service: service,
	}
}

func (r Resource) FindOne(oid string, req api2go.Request) (api2go.Responder, error) {
	id, err := domain.IDFrom(oid)
	if err != nil {
		return nil, resource.WrapError(err)
	}

	payment, err := r.service.Load(req.PlainRequest.Context(), id)
	if err != nil {
		return nil, resource.WrapError(err)
	}

	return resource.WrapObject(payment, http.StatusOK), nil
}

func (r Resource) FindAll(req api2go.Request) (api2go.Responder, error) {
	filter, err := r.ExtractSearchFilter(req.QueryParams)
	if err != nil {
		return nil, resource.WrapError(err)
	}

	searchResp, err := r.service.Search(req.PlainRequest.Context(), domain.PaymentSearchRequest{
		SearchFilter: filter,
	})
	if err != nil {
		return nil, resource.WrapError(err)
	}

	return resource.WrapArray(searchResp.Data, http.StatusOK), nil
}

func (r Resource) PaginatedFindAll(req api2go.Request) (uint, api2go.Responder, error) {
	filter, err := r.ExtractSearchFilter(req.QueryParams)
	if err != nil {
		return 0, nil, err
	}
	pagination, err := r.ExtractPagination(req.Pagination)
	if err != nil {
		return 0, nil, err
	}

	searchResp, err := r.service.Search(req.PlainRequest.Context(), domain.PaymentSearchRequest{
		SearchFilter:     filter,
		SearchPagination: pagination,
	})
	if err != nil {
		return 0, nil, resource.WrapError(err)
	}

	return searchResp.Size, resource.WrapArray(searchResp.Data, http.StatusOK), nil
}

func (r Resource) Create(obj interface{}, req api2go.Request) (api2go.Responder, error) {
	payment := obj.(*domain.Payment)

	err := payment.Validate()
	if err != nil {
		return nil, resource.WrapError(err)
	}

	err = r.service.Create(req.PlainRequest.Context(), payment)
	if err != nil {
		return nil, resource.WrapError(err)
	}

	return resource.WrapObject(payment, http.StatusCreated), nil
}

func (r Resource) Delete(oid string, req api2go.Request) (api2go.Responder, error) {
	id, err := domain.IDFrom(oid)
	if err != nil {
		return nil, resource.WrapError(err)
	}

	err = r.service.Delete(req.PlainRequest.Context(), id)
	if err != nil {
		return nil, resource.WrapError(err)
	}

	return resource.WrapObject(nil, http.StatusNoContent), nil
}

func (r Resource) Update(obj interface{}, req api2go.Request) (api2go.Responder, error) {
	payment := obj.(*domain.Payment)

	err := r.service.Update(req.PlainRequest.Context(), payment)
	if err != nil {
		return nil, resource.WrapError(err)
	}

	return resource.WrapObject(payment, http.StatusOK), nil
}

func paymentParamFunc(key string, values []string) (interface{}, error) {
	switch key {
	case "id":
		var ids []domain.ID
		for i, s := range values {
			id, err := domain.IDFrom(s)
			if err != nil {
				return nil, errors.Generic(
					errors.ErrCodeGenericInvalidArgument,
					err.Error(),
					fmt.Sprintf("field %q: index %d: %q has not invalid format", key, i, s),
				)
			}
			ids = append(ids, id)
		}
		return ids, nil
	case "creditor.account_number",
		"debtor.account_number":
		return values, nil
	default:
		return nil, errors.Generic(
			errors.ErrCodeGenericInvalidArgument,
			"unsupported filter parameter",
			key,
		)
	}
}
