package utils

import "errors"

type HTTPError struct {
	error

	errorCode uint32
}

func NewHTTPError(errorCode uint32, err string) *HTTPError {
	return &HTTPError{
		error:     errors.New(err),
		errorCode: errorCode,
	}
}

func (e HTTPError) Unwrap() error {
	return e.error
}

func WrapHTTPError(err error) *HTTPError {
	return WrapHTTPErrorWithCode(500, err)
}

func WrapHTTPErrorWithCode(errorCode uint32, err error) *HTTPError {
	return &HTTPError{
		error:     err,
		errorCode: errorCode,
	}
}
