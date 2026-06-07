package errors

import "net/http"

type Kind string

const (
	Validation   Kind = "validation"
	Unauthorized Kind = "unauthorized"
	Forbidden    Kind = "forbidden"
	NotFound     Kind = "not_found"
	Conflict     Kind = "conflict"
	Internal     Kind = "internal"
)

type Error struct {
	Kind    Kind
	Code    int
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return string(e.Kind)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func New(kind Kind, code int, message string, err error) *Error {
	return &Error{Kind: kind, Code: code, Message: message, Err: err}
}

func WithMessage(kind Kind, message string) *Error {
	return New(kind, 0, message, nil)
}

func HTTPStatus(kind Kind) int {
	switch kind {
	case Validation:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	case Forbidden:
		return http.StatusForbidden
	case NotFound:
		return http.StatusNotFound
	case Conflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
