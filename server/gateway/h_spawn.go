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
		rslt.JoinFmap(result.ParseSimplifiedNode),
		rslt.Fmap(func(r spawnRequest) string { return r.Object }),
		httpsrv.RequestJsonBody[spawnRequest],
	)

	return httpsrv.CurriedHandler(func(request *http.Request) func(writer http.ResponseWriter) {
		node, err := requestToSyntaxTree(request).Return()
		if err != nil {
			return httpsrv.RespondErrorWriter(httpsrv.BadRequest(err))
		}

		if err = objRunner.Run(request.Context(), node); err != nil {
			return httpsrv.RespondErrorWriter(err)
		}
		return httpsrv.RespondMessageWriter(http.StatusCreated, "spawned")
	})
}
