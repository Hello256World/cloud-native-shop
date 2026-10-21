package apierrors

type wrapped struct {
	APIError
	cause error
}

func (w *wrapped) Error() string { return w.cause.Error() }
func (w *wrapped) Unwrap() error { return w.cause }

func Wrap(api APIError, cause error) APIError {
	return &wrapped{APIError: api, cause: cause}
}
