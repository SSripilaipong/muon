package eventsource

type Observer interface {
	Update(events []AppendedEvent)
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

func (o *observeSubject) Update(events []AppendedEvent) {
	for _, observer := range o.observers {
		go observer.Update(events)
	}
}
