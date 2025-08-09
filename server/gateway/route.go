package gateway

import (
	"net/http"

	"github.com/SSripilaipong/muon/server/runner"
)

func newRouter(objRunner runner.Runner) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("POST /objects/spawn", spawnHandler(objRunner))
	return mux
}
