package handlers

import (
	"errors"
	"net/http"

	"github.com/Simaky/go-github-tracker/backend/app"
)

// Client-observable error codes. These are stable and may be relied on by API
// consumers.
const (
	CodeInvalidRequest = "INVALID_REQUEST"
	CodeNotFound       = "NOT_FOUND"
	CodeInternalError  = "INTERNAL_ERROR"
)

// FieldError describes a single field-level validation problem.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// APIError is the internal representation of a client-facing error. The cause
// is logged, never serialised.
type APIError struct {
	Status  int
	Code    string
	Message string
	Details []FieldError
	Cause   error
}

func (e *APIError) Error() string { return e.Message }
func (e *APIError) Unwrap() error { return e.Cause }

// BadRequest builds a 400 envelope.
func BadRequest(code, msg string, details ...FieldError) *APIError {
	return &APIError{Status: http.StatusBadRequest, Code: code, Message: msg, Details: details}
}

// NotFound builds a 404 envelope.
func NotFound(code, msg string) *APIError {
	return &APIError{Status: http.StatusNotFound, Code: code, Message: msg}
}

// InternalError builds a 500 envelope, retaining the cause for logging.
func InternalError(cause error) *APIError {
	return &APIError{Status: http.StatusInternalServerError, Code: CodeInternalError, Message: "internal error", Cause: cause}
}

// errorBody is the wire envelope the client receives.
type errorBody struct {
	Error struct {
		Code    string       `json:"code"`
		Message string       `json:"message"`
		Details []FieldError `json:"details,omitempty"`
	} `json:"error"`
}

// asAPIError maps any error to an *APIError, translating domain errors to the
// right status code.
func asAPIError(err error) *APIError {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}

	var domainErr *app.DomainError
	if errors.As(err, &domainErr) {
		switch domainErr.Code {
		case app.CodeValidation:
			return BadRequest(CodeInvalidRequest, domainErr.Message)
		case app.CodeNotFound:
			return NotFound(CodeNotFound, domainErr.Message)
		default:
			return InternalError(err)
		}
	}
	return InternalError(err)
}
