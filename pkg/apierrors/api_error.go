package apierrors

import "errors"

type Kind int

const (
	KindInternal Kind = iota
	KindInvalidInput
	KindNotFound
	KindConflict
	KindUnauthorized
	KindForbidden
	KindUnimplemented
)

type FieldError struct {
	Field string `json:"field"`
	Code  string `json:"code"`
	Param string `json:"param,omitempty"`
}

type APIError interface {
	error
	Kind() Kind
	Code() string
	Message() string
	Details() []FieldError
}

type apiError struct {
	kind    Kind
	code    string
	message string
	details []FieldError
}

func (e *apiError) Error() string   { return e.message }
func (e *apiError) Kind() Kind      { return e.kind }
func (e *apiError) Code() string    { return e.code }
func (e *apiError) Message() string { return e.message }
func (e *apiError) Details() []FieldError { return e.details }

func (e *apiError) Is(target error) bool {
	var t *apiError
	if errors.As(target, &t) {
		return e.code == t.code
	}
	return false
}
