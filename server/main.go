package server

import (
	"fmt"
	"log"

	"github.com/SSripilaipong/muon/common/system"
	"github.com/SSripilaipong/muon/server/coordinator"
	"github.com/SSripilaipong/muon/server/eventsource"
	"github.com/SSripilaipong/muon/server/gateway"
	"github.com/SSripilaipong/muon/server/runner"
)

func Start() error {
	esCtrl := eventsource.New()
	coordCtrl := coordinator.New(esCtrl)
	orCtrl := runner.New(esCtrl, coordCtrl)
	gw := gateway.New(runner.NewService(orCtrl))

	err, stopCoord := startCoordinator(coordCtrl)
	if err != nil {
		return err
	}
	defer stopCoord()

	err, stopEs := startEventSource(esCtrl)
	if err != nil {
		return err
	}
	defer stopEs()

	err, stopGateway := startGateway(gw)
	if err != nil {
		return err
	}
	defer stopGateway()

	err, stopRunner := startRunner(orCtrl)
	if err != nil {
		return err
	}
	defer stopRunner()

	select {
	case <-system.WaitForInterrupt():
	case <-esCtrl.Done():
	case <-gw.Done():
	case <-orCtrl.Done():
	}

	return nil
}

func startCoordinator(coord *coordinator.Controller) (error, func()) {
	if err := coord.Start(); err != nil {
		return fmt.Errorf("cannot start coordinator: %w", err), nil
	}
	return nil, func() {
		if err := coord.Stop(); err != nil {
			log.Println("stopping coordinator failed:", err)
		}
	}
}

func startEventSource(eventSource *eventsource.Controller) (error, func()) {
	if err := eventSource.Start(); err != nil {
		return fmt.Errorf("cannot start event source: %w", err), nil
	}
	return nil, func() {
		if err := eventSource.Stop(); err != nil {
			log.Println("stopping event source failed:", err)
		}
	}
}

func startRunner(objRunner *runner.Controller) (error, func()) {
	if err := objRunner.Start(); err != nil {
		return fmt.Errorf("cannot start object runner: %w", err), nil
	}
	return nil, func() {
		if err := objRunner.Stop(); err != nil {
			log.Println("stopping object runner:", err)
		}
	}
}

func startGateway(gw *gateway.Gateway) (error, func()) {
	if err := gw.Start(); err != nil {
		return fmt.Errorf("cannot start api gateway: %w", err), nil
	}
	return nil, func() {
		if err := gw.Stop(); err != nil {
			log.Println("stopping api gateway failed:", err)
		}
	}
}
