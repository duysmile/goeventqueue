# go-pubsub

Internal event queue with pub/sub pattern in Go with goroutines and channels

### Usage

Create queue to communicate between publisher and subscriber
```go
q := eventqueue.NewLocalQueue()
```


Create publisher to push event to queue
```go
pub := publisher.NewPublisher(q)
```
	
Create subscriber to consume event
```go
// `MaxGoRoutine` to control max routines to consumer event
// `MaxRetry` to control max backoff times to re-consumer event if it failed
sub := subscriber.NewSubscriber(q, subscriber.Config{
    MaxGoRoutine: 2,
    MaxRetry:     0,
})
```

Add handler to according event
```go
sub.Register(TestEvent, func(ctx context.Context, data interface{}) error {
    log.Println("job 1", data)
    wg.Done()
    return nil
})
```

Run workers to consume event
```go
sub.Start(mainCtx)
```

Push event to queue
```go
pub.Publish(mainCtx, NewEvent(TestEvent, "say"))
```


