package gateway

import (
	"net/http"
)

func NewRouter(objRunner Runner) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("POST /objects/spawn", spawnHandler(objRunner))
	return mux
}
