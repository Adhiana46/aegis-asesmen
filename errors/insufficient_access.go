package errors

import "net/http"

type ErrorInsufficientAccess struct {
	msg string
}

func NewErrorInsufficientAccess(msg ...string) error {
	e := ErrorInsufficientAccess{}

	if len(msg) > 0 {
		e.msg = msg[0]
	}

	return &e
}

func (e *ErrorInsufficientAccess) Error() string {
	if e.msg != "" {
		return e.msg
	}

	return "insufficient access"
}

func (e *ErrorInsufficientAccess) HttpStatusCode() int {
	return http.StatusBadRequest
}
