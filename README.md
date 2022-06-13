# Go Event Queue

Internal event queue with pub/sub pattern in Go with goroutines and channels

### Usage

Create queue to communicate between publisher and subscriber
```go
q := goeventqueue.NewLocalQueue(2)
```


Create publisher to push event to queue
```go
pub := goeventqueue.NewPublisher(q)
```
	
Create subscriber to consume event
```go
// `MaxGoRoutine` to control max routines to consumer event
// `MaxRetry` to control max backoff times to re-consumer event if it failed
sub := goeventqueue.NewSubscriber(q, subscriber.Config{
    MaxGoRoutine: 2,
    MaxRetry:     0,
})
```

Add handler to according event
```go
sub.Register(TestEvent, func(ctx context.Context, data interface{}) error {
    log.Println("job 1", data)
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

## License
MIT

## Contribution
All your contributions to project and make it better, they are welcome. Feel free to start an [issue](https://github.com/duysmile/goeventqueue/issues).

