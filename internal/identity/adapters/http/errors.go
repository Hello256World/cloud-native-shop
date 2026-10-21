package http

import (
	"errors"

	"github.com/Hello256World/cloud-native-shop/internal/identity/domain"
	"github.com/Hello256World/cloud-native-shop/pkg/apierrors"
)

func toAPIError(err error) error {
	var api apierrors.APIError

	switch {
	case errors.Is(err, domain.ErrPhoneTaken):
		api = apierrors.NewConflictError("phone_taken", "phone already registered")
	case errors.Is(err, domain.ErrEmailTaken):
		api = apierrors.NewConflictError("email_taken", "email already registered")
	case errors.Is(err, domain.ErrInvalidCredentials):
		api = apierrors.NewUnauthorizedError("invalid_credentials", "invalid phone or password")
	case errors.Is(err, domain.ErrRefreshTokenInvalid):
		api = apierrors.NewUnauthorizedError("refresh_token_invalid", "refresh token invalid or expired")
	case errors.Is(err, domain.ErrUserInactive):
		api = apierrors.NewForbiddenError("user_inactive", "user is inactive")
	case errors.Is(err, domain.ErrUserNotFound):
		api = apierrors.NewNotFoundError("user_not_found", "user not found")
	case errors.Is(err, domain.ErrAddressNotFound):
		api = apierrors.NewNotFoundError("address_not_found", "address not found")
	default:
		return err
	}

	return apierrors.Wrap(api, err)
}
