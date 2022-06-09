package eventqueue

type Queue interface {
	GetEventChan() chan Event
}

type localQueue struct {
	events chan Event
}

func (q *localQueue) GetEventChan() chan Event {
	return q.events
}

func NewLocalQueue() Queue {
	return &localQueue{
		events: make(chan Event),
	}
}
