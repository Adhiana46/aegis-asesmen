package errors

import "net/http"

type InternalError interface {
	Error() string
	HttpStatusCode() int
}

type InternalErrorStruct struct{}

func NewInternalError() InternalError {
	return &InternalErrorStruct{}
}

func (e *InternalErrorStruct) Error() string {
	return "internal error"
}

func (e *InternalErrorStruct) HttpStatusCode() int {
	return http.StatusInternalServerError
}
