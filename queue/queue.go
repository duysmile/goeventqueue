package queue

import "github.com/duysmile/goeventqueue"

// Queue exposes a channel for delivering events to subscribers.
type Queue interface {
        GetEventChan() chan goeventqueue.Event
}

type localQueue struct {
	events chan goeventqueue.Event
}

func (q *localQueue) GetEventChan() chan goeventqueue.Event {
	return q.events
}

// NewLocalQueue creates an in-memory queue with the given buffer size.
func NewLocalQueue(size int) Queue {
	return &localQueue{
		events: make(chan goeventqueue.Event, size),
	}
}
