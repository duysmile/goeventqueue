package goeventqueue

// Queue provides access to an underlying channel used to deliver events to
// subscribers.
type Queue interface {
	GetEventChan() chan Event
}

type localQueue struct {
	events chan Event
}

func (q *localQueue) GetEventChan() chan Event {
	return q.events
}

// NewLocalQueue creates a new in-memory queue with the provided buffer size.
// The queue implementation is safe for concurrent use by multiple publishers
// and subscribers.
func NewLocalQueue(size int) Queue {
	return &localQueue{
		events: make(chan Event, size),
	}
}
