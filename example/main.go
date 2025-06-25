package main

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/duysmile/goeventqueue"
	"github.com/duysmile/goeventqueue/publisher"
	"github.com/duysmile/goeventqueue/queue"
	"github.com/duysmile/goeventqueue/subscriber"
	"go.uber.org/zap"
)

type testEvent struct {
	Name goeventqueue.EventName
	Data string
}

func (t testEvent) GetName() goeventqueue.EventName {
	return t.Name
}

func (t testEvent) GetData() interface{} {
	return t.Data
}

func NewEvent(name goeventqueue.EventName, data string) goeventqueue.Event {
	return &testEvent{
		Name: name,
		Data: data,
	}
}

type logger struct {
	l *zap.Logger
}

func (l logger) Error(msg string, err error) {
	l.l.Error(msg, zap.Error(err))
}

func main() {
	var TestEvent goeventqueue.EventName = "test-event"

	q := queue.NewLocalQueue(2)

	pub := publisher.NewPublisher(q)
	sub := subscriber.NewSubscriber(q, subscriber.Config{
		MaxGoRoutine: 2,
		MaxRetry:     1,
	})

	zLog, _ := zap.NewDevelopment()
	logger := &logger{
		l: zLog,
	}
	sub.WithLogger(logger)

	mainCtx := context.Background()

	wg := sync.WaitGroup{}

	wg.Add(3)
	sub.Register(TestEvent, func(ctx context.Context, data interface{}) error {
		validData, ok := data.(string)
		if !ok {
			return nil
		}
		log.Println("job 1", validData)
		wg.Done()
		return nil
	})
	sub.Register(TestEvent, func(ctx context.Context, data interface{}) error {
		validData, ok := data.(string)
		if !ok {
			return nil
		}
		log.Println("job 2", validData)
		wg.Done()
		return nil
	})
	sub.Register(TestEvent, func(ctx context.Context, data interface{}) error {
		validData, ok := data.(string)
		if !ok {
			return nil
		}
		log.Println("job 3", validData)
		time.Sleep(2 * time.Second)
		wg.Done()
		return errors.New("error")
	})

	sub.Start(mainCtx)
	defer sub.Stop()

	_ = pub.Publish(mainCtx, NewEvent(TestEvent, "hello"))
	wg.Wait()
	time.Sleep(2 * time.Second)
}
