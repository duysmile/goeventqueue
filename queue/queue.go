package queue

import "github.com/duysmile/goeventqueue"

type Queue interface {
	GetEventChan() chan goeventqueue.Event
}

type localQueue struct {
	events chan goeventqueue.Event
}

func (q *localQueue) GetEventChan() chan goeventqueue.Event {
	return q.events
}

func NewLocalQueue(size int) Queue {
	return &localQueue{
		events: make(chan goeventqueue.Event, size),
	}
}
