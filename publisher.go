package goeventqueue

import (
	"context"
)

// Publisher posts events to a Queue so that they may be processed by
// subscribers.
type Publisher interface {
	Publish(ctx context.Context, event Event) error
}

type publisher struct {
	queue Queue
}

func (p publisher) Publish(ctx context.Context, event Event) error {
	go func(event Event) {
		select {
		case <-ctx.Done():
		case p.queue.GetEventChan() <- event:
		}
	}(event)

	return nil
}

// NewPublisher returns a Publisher that writes events to the provided Queue.
func NewPublisher(q Queue) Publisher {
	return &publisher{
		queue: q,
	}
}
