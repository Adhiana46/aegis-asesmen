package errors

import "net/http"

type ErrorInvalidCredentials struct {
	msg string
}

func NewErrorInvalidCredentials(msg ...string) error {
	e := ErrorInvalidCredentials{}

	if len(msg) > 0 {
		e.msg = msg[0]
	}

	return &e
}

func (e *ErrorInvalidCredentials) Error() string {
	if e.msg != "" {
		return e.msg
	}

	return "invalid credentials"
}

func (e *ErrorInvalidCredentials) HttpStatusCode() int {
	return http.StatusBadRequest
}
