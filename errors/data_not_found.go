package errors

import "net/http"

type ErrorDataNotFound struct {
	msg string
}

func NewErrorDataNotFound(msg ...string) error {
	e := ErrorDataNotFound{}

	if len(msg) > 0 {
		e.msg = msg[0]
	}

	return &e
}

func (e *ErrorDataNotFound) Error() string {
	if e.msg != "" {
		return e.msg
	}

	return "data not found"
}

func (e *ErrorDataNotFound) HttpStatusCode() int {
	return http.StatusBadRequest
}
