package httpsrv

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/fn"
	"github.com/SSripilaipong/muon/common/ioutil"
	"github.com/SSripilaipong/muon/common/jsonutil"
)

func RequestBody(request *http.Request) rslt.Of[io.ReadCloser] {
	if request.Body == nil {
		return rslt.Error[io.ReadCloser](fmt.Errorf("request.Body is nil"))
	}
	return rslt.Value(request.Body)
}

func RequestJsonBody[T any](request *http.Request) rslt.Of[T] {
	request.Context()
	return fn.Compose3(
		ResultWithBadRequest(rslt.JoinFmap(jsonutil.Read[T])),
		rslt.Fmap(ioutil.ToReader[io.ReadCloser]),
		RequestBody,
	)(request)
}

func RequestContext(request *http.Request) context.Context {
	return request.Context()
}
