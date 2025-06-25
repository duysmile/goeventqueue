package publisher

import (
	"context"

	"github.com/duysmile/goeventqueue"
	"github.com/duysmile/goeventqueue/queue"
)

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

func NewPublisher(q queue.Queue) Publisher {
	return &publisher{
		queue: q,
	}
}
