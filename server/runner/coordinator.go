package runner

import (
	"context"
	"fmt"
	"log"

	"github.com/SSripilaipong/muon/common/chn"
	"github.com/SSripilaipong/muon/common/prl"
)

type RemoteNode interface {
	ForceAppend(ctx context.Context, actions []Action) error
}

type Coordinator struct {
	local   *eventLog
	remotes []RemoteNode
}

type SubmitResult struct {
	LatestCommittedSequence uint64
	ChainedEvents           []ChainedEvent
	CommittedEvents         []AppendedEvent
}

func newCoordinator(local *eventLog, remotes []RemoteNode) *Coordinator {
	return &Coordinator{
		local:   local,
		remotes: remotes,
	}
}

func (c *Coordinator) Submit(ctx context.Context, actions []Action) (SubmitResult, error) {
	appendResp, appendErr := c.local.Append(actions)
	if appendErr != nil {
		return SubmitResult{}, fmt.Errorf("cannot append to local node: %w", appendErr)
	}

	remoteResponses := prl.Collect(ctx, appendingResultFromNodes(ctx, c.remotes, actions)...)
	ok, errs := guaranteeAtLeastHalfSuccess(ctx, c.numberOfOnlineNodes(), remoteResponses)
	if !ok {
		log.Println("[runner.coordinator] appending to other nodes failed")
		return SubmitResult{}, fmt.Errorf("quorum not satisfied: %v", errs)
	}

	committedEvents, commitErr := c.local.CommitUntil(appendResp.LatestCommittedSequence)
	if commitErr != nil {
		log.Println("[runner.coordinator] fail to mark commit until on local node")
		return SubmitResult{}, fmt.Errorf("cannot mark commit: %w", commitErr)
	}

	return SubmitResult{
		LatestCommittedSequence: appendResp.LatestCommittedSequence,
		ChainedEvents:           appendResp.ChainedEvents,
		CommittedEvents:         committedEvents,
	}, nil
}

func (c *Coordinator) numberOfOnlineNodes() int {
	return len(c.remotes) + 1
}

func appendingResultFromNodes(ctx context.Context, nodes []RemoteNode, actions []Action) (result []func() error) {
	for _, node := range nodes {
		node := node
		result = append(result, func() error {
			return node.ForceAppend(ctx, actions)
		})
	}
	return result
}

func guaranteeAtLeastHalfSuccess(ctx context.Context, n int, sources <-chan error) (bool, []error) {
	defer func() { go func() { chn.Drain(sources) }() }()

	var errs []error
loop:
	for {
		select {
		case err, isOpen := <-sources:
			if !isOpen {
				break loop
			}

			if err != nil {
				errs = append(errs, err)
			}
			if float64(len(errs)) >= float64(n)/2 {
				return false, errs
			}
		case <-ctx.Done():
			return false, errs
		}
	}
	return true, errs
}
