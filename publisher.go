package goeventqueue

import (
	"context"
)

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

func NewPublisher(q Queue) Publisher {
	return &publisher{
		queue: q,
	}
}
