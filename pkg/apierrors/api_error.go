package apierrors

type APIError interface {
	Error() string
	Message() string
	StatusCode() int
}

type apiError struct {
	XStatusCode int    `json:"-"`
	XMessage    string `json:"message"`
}

func (e apiError) Error() string { return e.XMessage }

func (e apiError) Message() string { return e.XMessage }

func (e apiError) StatusCode() int { return e.XStatusCode }
