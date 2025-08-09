package httpsrv

import (
	"net/http"

	"github.com/SSripilaipong/go-common/rslt"
)

func ResultWithBadRequest[A, B any](f func(rslt.Of[A]) rslt.Of[B]) func(rslt.Of[A]) rslt.Of[B] {
	return WrapError[A, B](BadRequest)(f)
}

func BadRequest(err error) error {
	return NewErrorWithStatusCode(http.StatusBadRequest, err)
}

func WrapError[A, B any](w func(error) error) func(func(of rslt.Of[A]) rslt.Of[B]) func(rslt.Of[A]) rslt.Of[B] {
	return func(f func(rslt.Of[A]) rslt.Of[B]) func(rslt.Of[A]) rslt.Of[B] {
		return func(x rslt.Of[A]) rslt.Of[B] {
			isErrorBefore := x.IsErr()
			y := f(x)
			isErrorAfter := y.IsErr()

			if !isErrorBefore && isErrorAfter {
				return rslt.Error[B](w(y.Error()))
			}
			return y
		}
	}
}

type ErrorWithStatusCode struct {
	statusCode int
	err        error
}

func (e ErrorWithStatusCode) Error() string {
	return e.err.Error()
}

func (e ErrorWithStatusCode) StatusCode() int {
	return e.statusCode
}

func NewErrorWithStatusCode(statusCode int, err error) ErrorWithStatusCode {
	return ErrorWithStatusCode{
		statusCode: statusCode,
		err:        err,
	}
}

func ErrorStatusCode(err error) int {
	x, hasStatusCode := err.(HasStatusCode)
	if !hasStatusCode {
		return http.StatusInternalServerError
	}
	return x.StatusCode()
}

type HasStatusCode interface {
	StatusCode() int
}
