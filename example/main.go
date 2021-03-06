package main

import (
	"context"
	"errors"
	"github.com/duysmile/goeventqueue"
	"go.uber.org/zap"
	"log"
	"sync"
	"time"
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

	q := goeventqueue.NewLocalQueue(2)

	pub := goeventqueue.NewPublisher(q)
	sub := goeventqueue.NewSubscriber(q, goeventqueue.Config{
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

	wg.Add(4)
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

	_ = pub.Publish(mainCtx, NewEvent(TestEvent, "hello"))
	wg.Wait()
	time.Sleep(2 * time.Second)
}
