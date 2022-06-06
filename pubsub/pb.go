package pubsub

import "context"

type Topic string

type Handler func(msg *Message) error

type PubSub interface {
	Publish(ctx context.Context, topic Topic, message *Message) error
	Subscribe(ctx context.Context, topic Topic) (ch <-chan *Message, close func())
}
