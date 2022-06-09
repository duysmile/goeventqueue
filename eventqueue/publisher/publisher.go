package publisher

import (
	"context"
	"github.com/duysmile/go-pubsub/eventqueue"
)

type Publisher interface {
	Publish(ctx context.Context, event eventqueue.Event) error
}

type publisher struct {
	queue eventqueue.Queue
}

func (p publisher) Publish(ctx context.Context, event eventqueue.Event) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.queue.GetEventChan() <- event:
		return nil
	}
}

func NewPublisher(q eventqueue.Queue) Publisher {
	return &publisher{
		queue: q,
	}
}
