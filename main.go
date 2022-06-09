package main

import (
	"context"
	"github.com/duysmile/go-pubsub/eventqueue"
	"github.com/duysmile/go-pubsub/eventqueue/publisher"
	"github.com/duysmile/go-pubsub/eventqueue/subscriber"
	"log"
	"sync"
	"time"
)

type testEvent struct {
	Name eventqueue.EventName
	Data interface{}
}

func (t testEvent) GetName() eventqueue.EventName {
	return t.Name
}

func (t testEvent) GetData() interface{} {
	return t.Data
}

func NewEvent(name eventqueue.EventName, data interface{}) eventqueue.Event {
	return &testEvent{
		Name: name,
		Data: data,
	}
}

func main() {
	var TestEvent eventqueue.EventName = "test-event"

	q := eventqueue.NewLocalQueue()

	pub := publisher.NewPublisher(q)
	sub := subscriber.NewSubscriber(q, subscriber.Config{
		MaxGoRoutine: 2,
		MaxRetry:     0,
	})

	mainCtx := context.Background()

	wg := sync.WaitGroup{}

	wg.Add(9)
	sub.Register(TestEvent, func(ctx context.Context, data interface{}) error {
		log.Println("job 1", data)
		wg.Done()
		return nil
	})
	sub.Register(TestEvent, func(ctx context.Context, data interface{}) error {
		log.Println("job 2", data)
		wg.Done()
		return nil
	})
	sub.Register(TestEvent, func(ctx context.Context, data interface{}) error {
		log.Println("job 3", data)
		time.Sleep(2 * time.Second)
		wg.Done()
		return nil
	})

	sub.Start(mainCtx)

	_ = pub.Publish(mainCtx, NewEvent(TestEvent, "hello"))
	_ = pub.Publish(mainCtx, NewEvent(TestEvent, "hi"))
	_ = pub.Publish(mainCtx, NewEvent(TestEvent, "bye"))
	wg.Wait()
}
