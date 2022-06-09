package subscriber

import (
	"context"
	"github.com/duysmile/goeventqueue"
	"log"
	"sync"
)

type Subscriber interface {
	Start(ctx context.Context)
	Register(name goeventqueue.EventName, handler Handler)
}

type Config struct {
	MaxGoRoutine int64
	MaxRetry     int64
}

type Handler func(ctx context.Context, data interface{}) error

type subscriber struct {
	queue           goeventqueue.Queue
	pool            chan goeventqueue.Event
	mapEventHandler map[goeventqueue.EventName][]Handler
	config          Config
	locker          sync.Mutex
}

func (s *subscriber) Register(name goeventqueue.EventName, handler Handler) {
	s.locker.Lock()

	listHandler, ok := s.mapEventHandler[name]
	if !ok {
		listHandler = make([]Handler, 0)
	}
	s.mapEventHandler[name] = append(listHandler, handler)
	s.locker.Unlock()
}

func (s *subscriber) Start(ctx context.Context) {
	for i := int64(0); i < s.config.MaxGoRoutine; i++ {
		go s.startWorker(ctx)
	}

	go s.startMainLoop(ctx)
}

func (s *subscriber) startMainLoop(ctx context.Context) {
	defer Recover()
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-s.queue.GetEventChan():
			if !ok {
				return
			}

			s.pool <- ev
		}
	}
}

func (s *subscriber) startWorker(ctx context.Context) {
	defer Recover()
	for {
		select {
		case <-ctx.Done():
			return
		case ev, ok := <-s.pool:
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
					select {
					case <-ctx.Done():
					case errChan <- job.Run(ctx, ev.GetData()):
					}
				}(ctx, handler)
			}

			for i := 0; i < len(listHandler); i++ {
				if err := <-errChan; err != nil {
					log.Println("Error handle job", err)
				}
			}
		}
	}
}

func NewSubscriber(q goeventqueue.Queue, cfg Config) Subscriber {
	return &subscriber{
		queue:           q,
		pool:            make(chan goeventqueue.Event, cfg.MaxGoRoutine),
		mapEventHandler: make(map[goeventqueue.EventName][]Handler),
		config:          cfg,
	}
}

func Recover() {
	if err := recover(); err != nil {
		log.Println(err)
	}
}
