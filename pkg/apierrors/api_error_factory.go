package apierrors

import (
	"errors"
	"net/http"
)

func NewAPIError(statusCode int, msg string) APIError {
	return apiError{
		XStatusCode: statusCode,
		XMessage:    msg,
	}
}

func NewBadRequestError(message string) APIError { return NewAPIError(http.StatusBadRequest, message) }

func NewNotFoundError(message string) APIError { return NewAPIError(http.StatusNotFound, message) }

func NewInternalServerError(message string) APIError {
	return NewAPIError(http.StatusInternalServerError, message)
}

func NewUnauthorizedError(message string) APIError {
	return NewAPIError(http.StatusUnauthorized, message)
}

func NewForbiddenError(message string) APIError { return NewAPIError(http.StatusForbidden, message) }

func NewUnimplementedError(message string) APIError {
	return NewAPIError(http.StatusNotImplemented, message)
}

func FromError(err error) APIError {
	if err == nil {
		return nil
	}

	var apiErr APIError
	if !errors.As(err, &apiErr) {
		return NewInternalServerError(err.Error())
	}

	return apiErr
}