package eventsource

// Store keeps track of appended events and commit progression for the local node.
type Store struct {
	events      []AppendedEvent
	commitUntil uint64
	observer    *observeSubject
}

func New() *Store {
	return &Store{observer: newObserveSubject()}
}

func (s *Store) AddObserver(observer Observer) {
	s.observer.Attach(observer)
}

func (s *Store) latestSequence() uint64 {
	if len(s.events) == 0 {
		return 0
	}
	return s.events[len(s.events)-1].Sequence()
}

func (s *Store) lastHash() uint64 {
	if len(s.events) == 0 {
		return 0
	}
	return s.events[len(s.events)-1].Hash()
}
