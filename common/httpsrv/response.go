package httpsrv

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func RespondErrorWriter(err error) func(w http.ResponseWriter) {
	return RespondMessageWriter(ErrorStatusCode(err), err.Error())
}

func RespondMessageWriter(statusCode int, message string) func(w http.ResponseWriter) {
	return func(w http.ResponseWriter) {
		w.WriteHeader(statusCode)
		writeMessageBody(message, w)
	}
}

func writeMessageBody(message string, w http.ResponseWriter) {
	r, err := json.Marshal(map[string]string{
		"message": message,
	})
	if err != nil {
		log.Println(fmt.Errorf("error while trying to marshal error response: %w", err))
		return
	}
	_, err = w.Write(r)
	if err != nil {
		log.Println(fmt.Errorf("error while trying to write error response: %w", err))
		return
	}
}
