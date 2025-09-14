package httpsrv

import (
	"net/http"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/rsltutil"
)

func ResultWithBadRequest[A, B any](f func(rslt.Of[A]) rslt.Of[B]) func(rslt.Of[A]) rslt.Of[B] {
	return rsltutil.WrapErrorFn[A, B](BadRequest)(f)
}

func BadRequest(err error) error {
	return NewErrorWithStatusCode(http.StatusBadRequest, err)
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
