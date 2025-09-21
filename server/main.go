package server

import (
	"context"
	"fmt"
	"log"

	"github.com/SSripilaipong/go-common/rslt"

	"github.com/SSripilaipong/muon/common/system"
	"github.com/SSripilaipong/muon/server/coordinator"
	"github.com/SSripilaipong/muon/server/eventsource"
	"github.com/SSripilaipong/muon/server/gateway"
	"github.com/SSripilaipong/muon/server/runner"
)

type runnerLocalNodeProxy struct {
	runner *runner.Controller
}

func (p *runnerLocalNodeProxy) LocalAppend(ctx context.Context, actions []eventsource.Action) rslt.Of[eventsource.AppendResponse] {
	if p.runner == nil {
		return rslt.Error[eventsource.AppendResponse](fmt.Errorf("runner is not initialized"))
	}
	return p.runner.LocalAppend(ctx, actions)
}

func (p *runnerLocalNodeProxy) MarkCommitUntil(ctx context.Context, sequence uint64) error {
	if p.runner == nil {
		return fmt.Errorf("runner is not initialized")
	}
	return p.runner.MarkCommitUntil(ctx, sequence)
}

func Start() error {
	esStore := eventsource.New()
	localProxy := &runnerLocalNodeProxy{}
	coordCtrl := coordinator.New(localProxy)
	orCtrl := runner.New(esStore, coordCtrl)
	localProxy.runner = orCtrl
	gw := gateway.New(runner.NewService(orCtrl))

	err, stopCoord := startCoordinator(coordCtrl)
	if err != nil {
		return err
	}
	defer stopCoord()

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
