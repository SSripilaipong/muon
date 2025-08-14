package runner

import (
	"slices"
	"time"

	es "github.com/SSripilaipong/muon/server/eventsource"
)

const channelTimeout = 500 * time.Millisecond

var runnerEventNames = []es.EventName{es.EventNameRun}

func isRunnerEventName(eventName es.EventName) bool {
	return slices.Contains(runnerEventNames, eventName)
}
