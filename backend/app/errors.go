package app

// Stable, client-observable failure modes. The server layer maps a
// *DomainError to the HTTP error envelope.
const (
	CodeValidation = "VALIDATION"
	CodeNotFound   = "NOT_FOUND"
	CodeConflict   = "CONFLICT"
	CodeUpstream   = "UPSTREAM"
	CodeInternal   = "INTERNAL"
)

// DomainError is the single typed error the domain layer returns. It carries a
// stable code and a client-safe message; the underlying cause is logged, never
// returned to the client.
type DomainError struct {
	Code    string
	Message string
	cause   error
}

func (e *DomainError) Error() string { return e.Message }
func (e *DomainError) Unwrap() error { return e.cause }

// NewValidationError reports invalid input.
func NewValidationError(msg string) *DomainError {
	return &DomainError{Code: CodeValidation, Message: msg}
}

// NewNotFoundError reports a missing resource ("<what> not found").
func NewNotFoundError(what string) *DomainError {
	return &DomainError{Code: CodeNotFound, Message: what + " not found"}
}

// NewConflictError reports a uniqueness/state conflict (e.g. a duplicate).
func NewConflictError(msg string) *DomainError {
	return &DomainError{Code: CodeConflict, Message: msg}
}

// NewUpstreamError reports a failure talking to an external dependency. The
// client-safe message is kept generic; the cause is logged.
func NewUpstreamError(msg string, cause error) *DomainError {
	return &DomainError{Code: CodeUpstream, Message: msg, cause: cause}
}

// WrapInternal wraps an unexpected cause as an internal domain error.
func WrapInternal(cause error) *DomainError {
	return &DomainError{Code: CodeInternal, Message: "internal error", cause: cause}
}
