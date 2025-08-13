package gateway

import (
	"context"
	"log"
	"net/http"
)

type Gateway struct {
	server   *http.Server
	stopChan chan struct{}
}

func New(objRunner Runner) *Gateway {
	return &Gateway{
		server: &http.Server{
			Addr:    ":8888",
			Handler: NewRouter(objRunner),
		},
		stopChan: make(chan struct{}),
	}
}

func (gw *Gateway) Start() error {
	go func() {
		defer close(gw.stopChan)

		if err := gw.server.ListenAndServe(); err != nil {
			log.Println("gateway stopped:", err)
		}
	}()
	return nil
}

func (gw *Gateway) Stop() error {
	if err := gw.server.Shutdown(context.Background()); err != nil {
		log.Println("error occurs while shutting down gateway:", err)
	}
	<-gw.Done()
	return nil
}

func (gw *Gateway) Done() chan struct{} {
	return gw.stopChan
}
