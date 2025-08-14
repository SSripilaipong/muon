package runner

import (
	"log"

	"github.com/SSripilaipong/muon/common/chn"
	es "github.com/SSripilaipong/muon/server/eventsource"
)

type EventSourceObserver struct {
	ctrl *Controller
}

func newEventSourceObserver(ctrl *Controller) es.Observer {
	return EventSourceObserver{ctrl: ctrl}
}

func (ob EventSourceObserver) Update(events []es.CommittedEvent) {
	for _, event := range events {
		if isRunnerEventName(event.EventName()) {
			if err := chn.SendWithTimeout[any](ob.ctrl.Ch(), event, channelTimeout); err != nil {
				log.Printf("cannot send committed event to runner: %v\n", err)
			}
		}
	}
}
