package httpsrv

import "net/http"

func CurriedHandler(f func(request *http.Request) func(writer http.ResponseWriter)) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		f(request)(writer)
	})
}
