package jsonutil

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/SSripilaipong/go-common/rslt"
)

func Read[T any](r io.Reader) rslt.Of[T] {
	decoder := json.NewDecoder(r)
	var x T
	err := decoder.Decode(&x)
	if err != nil {
		return rslt.Error[T](fmt.Errorf("cannot decode json: %w", err))
	}
	return rslt.Value(x)
}
