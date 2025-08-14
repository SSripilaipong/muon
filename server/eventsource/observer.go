package eventsource

type Observer interface {
	Update(events []CommittedEvent)
}

type observeSubject struct {
	observers []Observer
}

func newObserveSubject() *observeSubject {
	return &observeSubject{}
}

func (o *observeSubject) Attach(observer Observer) {
	o.observers = append(o.observers, observer)
}

func (o *observeSubject) Update(events []CommittedEvent) {
	for _, observer := range o.observers {
		go observer.Update(events)
	}
}
