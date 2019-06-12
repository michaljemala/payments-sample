package resource

import (
	"regexp"
	"strconv"

	"github.com/michaljemala/payments-sample/pkg/internal/errors"
)

var (
	filterParamPattern     = regexp.MustCompile(`^filter\[([\w._]+)\]$`)
	paginationParamPattern = regexp.MustCompile(`^page\[([number|size]+)\]$`)
)

type Generic struct {
	ParamFunc func(key string, values []string) (interface{}, error)
}

func (r *Generic) ExtractSearchFilter(searchFilterParams map[string][]string) (SearchFilter, error) {
	filter := make(SearchFilter)

	for key, values := range searchFilterParams {
		tokens := paginationParamPattern.FindStringSubmatch(key)
		if len(tokens) == 2 {
			continue
		}

		tokens = filterParamPattern.FindStringSubmatch(key)
		if len(tokens) != 2 {
			return nil, errors.Generic(
				errors.ErrCodeGenericInvalidArgument,
				"invalid query parameter",
				key,
			)
		}
		key = tokens[1]

		extracted, err := r.ParamFunc(key, values)
		if err != nil {
			return nil, errors.Generic(
				errors.ErrCodeGenericInvalidArgument,
				"invalid filter parameter",
				key,
			)
		}

		filter[key] = extracted
	}

	return filter, nil
}

func (r *Generic) ExtractPagination(paginationParams map[string]string) (*SearchPagination, error) {
	var (
		page, size uint64
		err        error
	)
	for key, value := range paginationParams {
		switch key {
		case "number":
			page, err = strconv.ParseUint(value, 0, 0)
			if err != nil || page == 0 {
				return nil, errors.Generic(
					errors.ErrCodeGenericInvalidArgument,
					"invalid paging parameter",
					key,
				)
			}
		case "size":
			size, err = strconv.ParseUint(value, 0, 0)
			if err != nil || size == 0 {
				return nil, errors.Generic(
					errors.ErrCodeGenericInvalidArgument,
					"invalid paging parameter",
					key,
				)
			}
		default:
			return nil, errors.Generic(
				errors.ErrCodeGenericInvalidArgument,
				"unsupported paging parameter",
				key,
			)
		}
	}
	return NewSearchPagination(uint(page), uint(size)), nil
}
