package errs

import (
	"errors"
	"fmt"
)

// TODO: do i need another err code for generic bad requests?
const (
	ERRCONFLICT       = "conflict"
	ERRINTERNAL       = "internal"
	ERRBAD            = "bad"
	ERRUNPROCESSABLE  = "unprocessable"
	ERRNOTFOUND       = "not_found"
	ERRNOTIMPLEMENTED = "not_implemented"
	ERRNOTAUTHORIZED  = "not_authorized"
)

type Error struct {
	Code    string
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("coffeetanuki error: code=%s message=%s", e.Code, e.Message)
}

func Errorf(code string, format string, args ...any) *Error {
	return &Error{
		Code:    code,
		Message: fmt.Sprintf(format, args...),
	}
}

func ErrorCode(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	} else {
		return ERRINTERNAL
	}
}

func ErrorMessage(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Message
	} else {
		return "internal error"
	}
}
