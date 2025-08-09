package server

import (
	"fmt"
	"log"

	"github.com/SSripilaipong/muon/common/system"
	"github.com/SSripilaipong/muon/server/gateway"
	"github.com/SSripilaipong/muon/server/runner"
)

func Start() error {
	objRunner := runner.New()
	gw := gateway.New(objRunner)

	err, stopGateway := startGateway(gw)
	if err != nil {
		return err
	}
	defer stopGateway()

	err, stopRunner := startRunner(objRunner)
	if err != nil {
		return err
	}
	defer stopRunner()

	select {
	case <-system.WaitForInterrupt():
	case <-gw.Done():
	case <-objRunner.Done():
	}

	return nil
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
