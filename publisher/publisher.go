package publisher

import (
	"context"

	"github.com/duysmile/goeventqueue"
	"github.com/duysmile/goeventqueue/queue"
)

// Publisher posts events to a queue so that they can be consumed by subscribers.
type Publisher interface {
	Publish(ctx context.Context, event goeventqueue.Event) error
}

type publisher struct {
	queue queue.Queue
}

func (p publisher) Publish(ctx context.Context, event goeventqueue.Event) error {
	go func(event goeventqueue.Event) {
		select {
		case <-ctx.Done():
		case p.queue.GetEventChan() <- event:
		}
	}(event)

	return nil
}

// NewPublisher returns a Publisher that writes events to the provided queue.
func NewPublisher(q queue.Queue) Publisher {
	return &publisher{
		queue: q,
	}
}
