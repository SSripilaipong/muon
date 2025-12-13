package runner

import "fmt"

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

type eventLog struct {
	events      []AppendedEvent
	commitUntil uint64
}

func newEventLog() *eventLog {
	return &eventLog{}
}

func (l *eventLog) Append(actions []Action) (AppendResponse, error) {
	previousHash := uint64(0)
	if len(l.events) > 0 {
		previousHash = l.events[len(l.events)-1].Hash()
	}
	appendedEvents, appendErr := processAppendActions(actions, l.latestSequence())
	if appendErr != nil {
		return AppendResponse{}, appendErr
	}

	l.events = append(l.events, appendedEvents...)
	latestSequence := l.latestSequence()

	chainedEvents := chainEvents(appendedEvents, previousHash, latestSequence)
	return AppendResponse{
		LatestCommittedSequence: latestSequence,
		ChainedEvents:           chainedEvents,
	}, nil
}

func (l *eventLog) CommitUntil(seq uint64) ([]AppendedEvent, error) {
	previousCommitUntil := l.commitUntil
	if seq < l.commitUntil {
		return nil, fmt.Errorf("decreasing commit is not allowed")
	}
	l.commitUntil = seq

	commitStartIndex, err := seekToSequence(l.events, previousCommitUntil+1)
	if err != nil {
		return nil, err
	}
	return l.events[commitStartIndex:], nil
}

func (l *eventLog) latestSequence() uint64 {
	if len(l.events) == 0 {
		return 0
	}
	return l.events[len(l.events)-1].Sequence()
}

func chainEvents(appendedEvents []AppendedEvent, latestHash uint64, latestSequence uint64) []ChainedEvent {
	var chainedEvents []ChainedEvent
	previousHash, previousSequence := latestHash, latestSequence
	for _, event := range appendedEvents {
		chainedEvents = append(chainedEvents, ChainedEvent{
			Event:            event,
			PreviousSequence: previousSequence,
			PreviousHash:     previousHash,
		})
		previousHash, previousSequence = event.Hash(), event.Sequence()
	}
	return chainedEvents
}

func processAppendActions(actions []Action, previousSeq uint64) ([]AppendedEvent, error) {
	var eventsToAppend []AppendedEvent
	seq := previousSeq
	for _, action := range actions {
		switch action := action.(type) {
		case AppendAction:
			if action.requiredSequence.IsNotEmpty() && action.requiredSequence.Value() != seq {
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

func seekToSequence(events []AppendedEvent, seq uint64) (uint64, error) {
	for i := len(events) - 1; i >= 0; i-- {
		event := events[i]
		if event.Sequence() == seq {
			return uint64(i), nil
		} else if event.Sequence() < seq {
			return 0, fmt.Errorf("out of sequence")
		}
	}
	return 0, fmt.Errorf("sequence not found")
}
