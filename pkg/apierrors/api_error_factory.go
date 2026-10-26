package apierrors

import "errors"

func New(kind Kind, code, message string) APIError {
	return &apiError{kind: kind, code: code, message: message}
}

func NewBadRequestError(code, message string) APIError   { return New(KindInvalidInput, code, message) }
func NewNotFoundError(code, message string) APIError     { return New(KindNotFound, code, message) }
func NewConflictError(code, message string) APIError     { return New(KindConflict, code, message) }
func NewUnauthorizedError(code, message string) APIError { return New(KindUnauthorized, code, message) }
func NewForbiddenError(code, message string) APIError    { return New(KindForbidden, code, message) }
func NewUnimplementedError(code, message string) APIError {
	return New(KindUnimplemented, code, message)
}
func NewValidationError(details []FieldError) APIError {
	return &apiError{
		kind:    KindInvalidInput,
		code:    "validation_failed",
		message: "request validation failed",
		details: details,
	}
}

func NewInternalServerError() APIError {
	return New(KindInternal, "internal_error", "internal server error")
}

func FromError(err error) APIError {
	if err == nil {
		return nil
	}
	var apiErr APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	return NewInternalServerError()
}
