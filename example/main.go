package main

import (
	"context"
	"errors"
	"github.com/duysmile/goeventqueue"
	"log"
	"sync"
	"time"
)

type testEvent struct {
	Name goeventqueue.EventName
	Data interface{}
}

func (t testEvent) GetName() goeventqueue.EventName {
	return t.Name
}

func (t testEvent) GetData() interface{} {
	return t.Data
}

func NewEvent(name goeventqueue.EventName, data interface{}) goeventqueue.Event {
	return &testEvent{
		Name: name,
		Data: data,
	}
}

func main() {
	var TestEvent goeventqueue.EventName = "test-event"

	q := goeventqueue.NewLocalQueue(2)

	pub := goeventqueue.NewPublisher(q)
	sub := goeventqueue.NewSubscriber(q, goeventqueue.Config{
		MaxGoRoutine: 2,
		MaxRetry:     1,
	})

	mainCtx := context.Background()

	wg := sync.WaitGroup{}

	wg.Add(3)
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
		return errors.New("error")
	})

	sub.Start(mainCtx)

	_ = pub.Publish(mainCtx, NewEvent(TestEvent, "hello"))
	wg.Wait()
}
