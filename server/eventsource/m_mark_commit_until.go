package eventsource

import (
	"context"
	"fmt"
	"log"

	"github.com/SSripilaipong/go-common/rslt"
)

func (s *Store) MarkCommitUntil(_ context.Context, sequence uint64) error {
	if sequence < s.commitUntil {
		return fmt.Errorf("decreasing commit is not allowed")
	}
	if sequence == s.commitUntil {
		return nil
	}

	previousCommitUntil := s.commitUntil
	s.commitUntil = sequence

	commitStartIndex, err := seekToSequence(s.events, previousCommitUntil+1).Return()
	if err != nil {
		log.Println("[server.eventsource] error while seeking sequence:", err)
		return nil
	}

	s.observer.Update(s.events[commitStartIndex:])

	return nil
}

func seekToSequence(events []AppendedEvent, seq uint64) rslt.Of[uint64] {
	for i := len(events) - 1; i >= 0; i-- {
		event := events[i]
		if event.Sequence() == seq {
			return rslt.Value(uint64(i))
		} else if event.Sequence() < seq {
			return rslt.Error[uint64](fmt.Errorf("out of sequence"))
		}
	}
	return rslt.Error[uint64](fmt.Errorf("sequence not found"))
}
