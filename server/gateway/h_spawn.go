package gateway

import (
	"net/http"

	"github.com/SSripilaipong/go-common/rslt"
	"github.com/SSripilaipong/muto/parser/result"

	"github.com/SSripilaipong/muon/common/fn"
	"github.com/SSripilaipong/muon/common/httpsrv"
)

func spawnHandler(objRunner Runner) http.Handler {
	type spawnRequest struct {
		Object string `json:"object,required"`
	}

	requestToSyntaxTree := fn.Compose3(
		httpsrv.ResultWithBadRequest(rslt.JoinFmap(result.ParseSimplifiedNode)),
		rslt.Fmap(func(r spawnRequest) string { return r.Object }),
		httpsrv.RequestJsonBody[spawnRequest],
	)
	spawn := fn.Compose(
		rslt.Transform(objRunner.Run, fn.Id),
		requestToSyntaxTree,
	)

	return httpsrv.CurriedHandler(func(request *http.Request) func(writer http.ResponseWriter) {
		if err := spawn(request); err != nil {
			return httpsrv.RespondErrorWriter(err)
		}
		return httpsrv.RespondMessageWriter(http.StatusCreated, "spawned")
	})
}
