package servererrors

import "errors"

var (
	ErrorRecordNotFound = errors.New("record not found")
	ErrorInternal       = errors.New("internal server error")
)
