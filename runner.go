package main

import (
	"context"
	"github.com/duysmile/go-pubsub/pubsub"
	"log"
)

var (
	SendEmailTopic pubsub.Topic = "send-email"
)

type HandlerFunc func(msg *pubsub.Message) error

type consumerJob struct {
	Title   string
	Handler HandlerFunc
}

type PubSubRunner struct {
	ps pubsub.PubSub
}

func NewPubSubRunner(ps pubsub.PubSub) *PubSubRunner {
	return &PubSubRunner{
		ps: ps,
	}
}

func NewConsumerJob(title string, handler HandlerFunc) *consumerJob {
	return &consumerJob{
		Title:   title,
		Handler: handler,
	}
}

func (r *PubSubRunner) startSubscribeTopic(ctx context.Context, topic pubsub.Topic, job *consumerJob) {
	subscribeChan, _ := r.ps.Subscribe(ctx, topic)
	go func() {
		for {
			msg := <-subscribeChan
			_ = job.Handler(msg)
		}
	}()
}

func (r *PubSubRunner) Init(ctx context.Context) {
	r.startSubscribeTopic(
		ctx,
		SendEmailTopic,
		NewConsumerJob(string(SendEmailTopic), func(msg *pubsub.Message) error {
			log.Println("Handle message", msg)
			return nil
		}),
	)
}
