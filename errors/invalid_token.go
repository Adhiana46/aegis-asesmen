package errors

import "net/http"

type ErrorInvalidToken struct {
	msg string
}

func NewErrorInvalidToken(msg ...string) error {
	e := ErrorInvalidToken{}

	if len(msg) > 0 {
		e.msg = msg[0]
	}

	return &e
}

func (e *ErrorInvalidToken) Error() string {
	if e.msg != "" {
		return e.msg
	}

	return "invalid token"
}

func (e *ErrorInvalidToken) HttpStatusCode() int {
	return http.StatusBadRequest
}
