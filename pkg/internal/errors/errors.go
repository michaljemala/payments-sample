package errors

import "fmt"

type errorCategory string

const (
	ErrCategoryGeneric    = errorCategory("GENERIC")
	ErrCategoryDataAccess = errorCategory("DATA_ACCESS")
)

type errorCodeGeneric string

func (c errorCodeGeneric) code() string   { return string(c) }
func (c errorCodeGeneric) String() string { return c.code() }

const (
	ErrCodeGenericAlreadyExists   = errorCodeGeneric("ALREADY_EXISTS")
	ErrCodeGenericInvalidArgument = errorCodeGeneric("INVALID_ARGUMENT")
	ErrCodeGenericInternal        = errorCodeGeneric("INTERNAL")
	ErrCodeGenericNotFound        = errorCodeGeneric("NOT_FOUND")
)

type errorCodeDataAccess string

func (c errorCodeDataAccess) code() string   { return string(c) }
func (c errorCodeDataAccess) String() string { return c.code() }

const (
	ErrCodeDataAccessInsertFailed = errorCodeDataAccess("INSERT_FAILED")
	ErrCodeDataAccessSelectFailed = errorCodeDataAccess("SELECT_FAILED")
	ErrCodeDataAccessDeleteFailed = errorCodeDataAccess("DELETE_FAILED")
	ErrCodeDataAccessUpdateFailed = errorCodeDataAccess("UPDATE_FAILED")
)

type Error struct {
	Category errorCategory          `json:"category"`
	Code     Code                   `json:"code"`
	Message  string                 `json:"message"`
	Detail   string                 `json:"detail"`
	Extra    map[string]interface{} `json:"extra"`
}

func (e Error) Error() string {
	return fmt.Sprintf(e.Message)
}

type Code interface {
	code() string
	String() string
}

func Generic(code errorCodeGeneric, msg string, detail string, extra ...map[string]interface{}) Error {
	return Error{
		Category: ErrCategoryGeneric,
		Code:     code,
		Message:  msg,
		Detail:   detail,
		Extra:    mergeMaps(extra...),
	}
}

func DataAccess(code errorCodeDataAccess, msg string, detail string, extra ...map[string]interface{}) Error {
	return Error{
		Category: ErrCategoryDataAccess,
		Code:     code,
		Message:  msg,
		Detail:   detail,
		Extra:    mergeMaps(extra...),
	}
}

func mergeMaps(extras ...map[string]interface{}) map[string]interface{} {
	var merged map[string]interface{}
	if len(extras) != 0 {
		merged = make(map[string]interface{})
		for _, extra := range extras {
			for k, v := range extra {
				merged[k] = v
			}
		}
	}
	return merged
}
