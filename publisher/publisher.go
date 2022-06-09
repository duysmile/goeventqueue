package publisher

import (
	"context"
	"github.com/duysmile/goeventqueue"
)

type Publisher interface {
	Publish(ctx context.Context, event goeventqueue.Event) error
}

type publisher struct {
	queue goeventqueue.Queue
}

func (p publisher) Publish(ctx context.Context, event goeventqueue.Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.queue.GetEventChan() <- event:
		return nil
	}
}

func NewPublisher(q goeventqueue.Queue) Publisher {
	return &publisher{
		queue: q,
	}
}
