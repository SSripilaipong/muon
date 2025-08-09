package ioutil

import "io"

func ToReader[T io.Reader](x T) io.Reader {
	return x
}
