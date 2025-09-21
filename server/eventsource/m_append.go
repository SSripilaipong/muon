package eventsource

import (
	"context"
	"fmt"

	"github.com/SSripilaipong/go-common/rslt"
)

type Action any

type AppendResponse struct {
	LatestCommittedSequence uint64
	ChainedEvents           []ChainedEvent
}

type ChainedEvent struct {
	Event            AppendedEvent
	PreviousSequence uint64
	PreviousHash     uint64
}

func (s *Store) LocalAppend(_ context.Context, actions []Action) rslt.Of[AppendResponse] {
	previousSequence := s.latestSequence()
	previousHash := s.lastHash()

	eventsToAppend, appendErr := processAppendActions(actions, previousSequence)
	if appendErr != nil {
		return rslt.Error[AppendResponse](appendErr)
	}

	s.events = append(s.events, eventsToAppend...)

	latestSequence := s.latestSequence()
	chainedEvents := buildChainedEvents(eventsToAppend, previousSequence, previousHash)

	return rslt.Value(AppendResponse{
		LatestCommittedSequence: latestSequence,
		ChainedEvents:           chainedEvents,
	})
}

func (s *Store) ForceAppend(ctx context.Context, actions []Action) rslt.Of[AppendResponse] {
	return s.LocalAppend(ctx, actions)
}

func buildChainedEvents(appended []AppendedEvent, previousSequence, previousHash uint64) []ChainedEvent {
	if len(appended) == 0 {
		return nil
	}

	chained := make([]ChainedEvent, 0, len(appended))
	seq, hash := previousSequence, previousHash
	for _, event := range appended {
		chained = append(chained, ChainedEvent{
			Event:            event,
			PreviousSequence: seq,
			PreviousHash:     hash,
		})
		seq = event.Sequence()
		hash = event.Hash()
	}
	return chained
}

func processAppendActions(actions []Action, previousSeq uint64) ([]AppendedEvent, error) {
	var eventsToAppend []AppendedEvent
	seq := previousSeq
	for _, action := range actions {
		switch action := action.(type) {
		case AppendAction:
			rs := action.requiredSequence
			if rs.IsNotEmpty() && rs.Value() != seq {
				return nil, fmt.Errorf("sequence requirement violation")
			}
			seq++
			eventsToAppend = append(eventsToAppend, NewAppended(action.event, seq))
		default:
			return nil, fmt.Errorf("unknown action %T", action)
		}
	}
	return eventsToAppend, nil
}
