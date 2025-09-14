package eventsource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_processAppendActions(t *testing.T) {
	t.Run("should return appended events", func(t *testing.T) {
		resultEvents, err := processAppendActions([]Action{
			NewAppendAction(RunEvent{ModuleVersion: "newer", hash: 123}),
			NewAppendAction(RunEvent{ModuleVersion: "newest", hash: 456}),
		}, 2)

		assert.Nil(t, err)
		assert.Equal(t, AppendedEvent{
			event:    RunEvent{ModuleVersion: "newer", hash: 123},
			sequence: 3,
		}, resultEvents[0])
		assert.Equal(t, AppendedEvent{
			event:    RunEvent{ModuleVersion: "newest", hash: 456},
			sequence: 4,
		}, resultEvents[1])
	})

	t.Run("should validate required sequence", func(t *testing.T) {
		_, err := processAppendActions([]Action{
			NewAppendAction(RunEvent{}, AppendAtSequence(3)),
		}, 2)

		assert.Equal(t, "sequence requirement violation", err.Error())
	})
}
