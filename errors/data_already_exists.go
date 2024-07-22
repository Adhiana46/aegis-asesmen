package errors

import "net/http"

type ErrorDataAlreadyExists struct {
	msg string
}

func NewErrorDataAlreadyExists(msg ...string) error {
	e := ErrorDataAlreadyExists{}

	if len(msg) > 0 {
		e.msg = msg[0]
	}

	return &e
}

func (e *ErrorDataAlreadyExists) Error() string {
	if e.msg != "" {
		return e.msg
	}

	return "data already exists"
}

func (e *ErrorDataAlreadyExists) HttpStatusCode() int {
	return http.StatusBadRequest
}
