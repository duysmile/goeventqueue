package subscriber

import (
	"context"
	"fmt"
	"sync"

	"github.com/duysmile/goeventqueue"
	"github.com/duysmile/goeventqueue/queue"
)

// Subscriber consumes events from a Queue and dispatches them to registered
// handlers.
type Subscriber interface {
	Start(ctx context.Context)
	Register(name goeventqueue.EventName, handler Handler)
	WithLogger(l Logger)
	Stop()
}

// Config defines parameters that control how the subscriber dispatches work.
type Config struct {
	MaxGoRoutine int64
	MaxRetry     int64
}

// Handler processes event data received by the Subscriber.
type Handler func(ctx context.Context, data interface{}) error

type subscriber struct {
	queue           queue.Queue
	mapEventHandler map[goeventqueue.EventName][]Handler
	config          Config
	locker          sync.Mutex
	logger          Logger
	quit            chan struct{}
	wg              sync.WaitGroup
}

func (s *subscriber) Register(name goeventqueue.EventName, handler Handler) {
	s.locker.Lock()
	defer s.locker.Unlock()

	listHandler, ok := s.mapEventHandler[name]
	if !ok {
		listHandler = make([]Handler, 0)
	}
	s.mapEventHandler[name] = append(listHandler, handler)
}

func (s *subscriber) Stop() {
	close(s.quit)
	s.wg.Wait()
}

func (s *subscriber) Start(ctx context.Context) {
	eQueue := s.queue.GetEventChan()
	for i := int64(0); i < s.config.MaxGoRoutine; i++ {
		s.wg.Add(1)
		go s.startWorker(ctx, eQueue)
	}
}

func (s *subscriber) startWorker(ctx context.Context, eQueue chan goeventqueue.Event) {
	defer s.Recover()
	defer s.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.quit:
			return
		case ev, ok := <-eQueue:
			if !ok {
				return
			}

			s.locker.RLock()
			listHandler, ok := s.mapEventHandler[ev.GetName()]
			s.locker.RUnlock()
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
					case errChan <- job.Run(ctx, ev.GetData()):
					case <-s.quit:
					}
				}(ctx, handler)
			}

			for i := 0; i < len(listHandler); i++ {
				if err := <-errChan; err != nil {
					s.logger.Error("error after all retry times", err)
				}
			}
		}
	}
}

// NewSubscriber creates a Subscriber that reads from the provided Queue using
// the supplied configuration.
func NewSubscriber(q queue.Queue, cfg Config) Subscriber {
	return &subscriber{
		queue:           q,
		mapEventHandler: make(map[goeventqueue.EventName][]Handler),
		config:          cfg,
		logger:          NewDefaultLogger(),
		quit:            make(chan struct{}),
	}
}

func (s *subscriber) WithLogger(l Logger) {
	s.logger = l
}

func (s *subscriber) Recover() {
	if err := recover(); err != nil {
		s.logger.Error("panic error", fmt.Errorf("%v", err))
	}
}
