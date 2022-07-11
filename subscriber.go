package goeventqueue

import (
	"context"
	"log"
	"sync"
)

type Subscriber interface {
	Start(ctx context.Context)
	Register(name EventName, handler Handler)
}

type Config struct {
	MaxGoRoutine int64
	MaxRetry     int64
}

type Handler func(ctx context.Context, data interface{}) error

type subscriber struct {
	queue           Queue
	mapEventHandler map[EventName][]Handler
	config          Config
	locker          sync.Mutex
}

func (s *subscriber) Register(name EventName, handler Handler) {
	s.locker.Lock()

	listHandler, ok := s.mapEventHandler[name]
	if !ok {
		listHandler = make([]Handler, 0)
	}
	s.mapEventHandler[name] = append(listHandler, handler)
	s.locker.Unlock()
}

func (s *subscriber) Start(ctx context.Context) {
	eQueue := s.queue.GetEventChan()
	for i := int64(0); i < s.config.MaxGoRoutine; i++ {
		go func(ctx context.Context, eQueue chan Event) {
			s.startWorker(ctx, eQueue)
		}(ctx, eQueue)
	}
}

func (s *subscriber) startWorker(ctx context.Context, eQueue chan Event) {
	defer Recover()
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-eQueue:
			if !ok {
				return
			}

			listHandler, ok := s.mapEventHandler[ev.GetName()]
			if !ok {
				continue
			}

			errChan := make(chan error, len(listHandler))
			for _, handler := range listHandler {
				go func(ctx context.Context, handler Handler) {
					job := NewJob(handler, JobConfig{
						MaxBackOff: s.config.MaxRetry,
					})

					errChan <- job.Run(ctx, ev.GetData())
				}(ctx, handler)
			}

			for i := 0; i < len(listHandler); i++ {
				if err := <-errChan; err != nil {
					log.Println("error after all retry times", err)
				}
			}
		}
	}
}

func NewSubscriber(q Queue, cfg Config) Subscriber {
	return &subscriber{
		queue:           q,
		mapEventHandler: make(map[EventName][]Handler),
		config:          cfg,
	}
}

func Recover() {
	if err := recover(); err != nil {
		log.Println(err)
	}
}
