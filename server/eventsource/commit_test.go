package eventsource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_processCommitActions(t *testing.T) {
	t.Run("should append events", func(t *testing.T) {
		existingEvents := []CommittedEvent{
			NewCommitted(RunEvent{ModuleVersion: "oldest"}, 1),
			NewCommitted(RunEvent{ModuleVersion: "older"}, 2),
		}

		resultEvents, commitSequences, err := processCommitActions([]Action{
			NewAppendAction(RunEvent{ModuleVersion: "newer"}),
			NewAppendAction(RunEvent{ModuleVersion: "newest"}),
		}, 2, existingEvents)

		assert.Nil(t, err)
		assert.Equal(t, []int64{3, 4}, commitSequences)
		assert.Equal(t, CommittedEvent{
			Event:    RunEvent{ModuleVersion: "newer"},
			sequence: 3,
		}, resultEvents[2])
		assert.Equal(t, CommittedEvent{
			Event:    RunEvent{ModuleVersion: "newest"},
			sequence: 4,
		}, resultEvents[3])
	})

	t.Run("should validate required sequence", func(t *testing.T) {
		existingEvents := []CommittedEvent{
			NewCommitted(RunEvent{}, 1),
			NewCommitted(RunEvent{}, 2),
		}

		_, _, err := processCommitActions([]Action{
			NewAppendAction(RunEvent{}, AppendAtSequence(3)),
		}, 2, existingEvents)

		assert.Equal(t, "sequence requirement violation", err.Error())
	})
}
