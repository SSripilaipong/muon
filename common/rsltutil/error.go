package rsltutil

import "github.com/SSripilaipong/go-common/rslt"

func WrapError[T any](w func(error) error) func(rslt.Of[T]) rslt.Of[T] {
	return func(x rslt.Of[T]) rslt.Of[T] {
		if x.IsErr() {
			return rslt.Error[T](w(x.Error()))
		}
		return x
	}
}

func WrapErrorFn[A, B any](w func(error) error) func(func(rslt.Of[A]) rslt.Of[B]) func(rslt.Of[A]) rslt.Of[B] {
	wrap := WrapError[B](w)
	return func(f func(rslt.Of[A]) rslt.Of[B]) func(rslt.Of[A]) rslt.Of[B] {
		return func(x rslt.Of[A]) rslt.Of[B] {
			isErrorBefore := x.IsErr()
			y := f(x)

			if !isErrorBefore {
				return wrap(y)
			}
			return y
		}
	}
}
