package plocal

import (
	"context"
	"github.com/duysmile/go-pubsub/pubsub"
	"log"
	"sync"
)

type localPubSub struct {
	messageQueue    chan *pubsub.Message
	mapTopicChannel map[pubsub.Topic][]chan *pubsub.Message
	locker          *sync.RWMutex
}

func (localP localPubSub) Publish(ctx context.Context, topic pubsub.Topic, message *pubsub.Message) error {
	message.SetTopic(topic)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Failed to publish message", message)
			}
		}()

		localP.messageQueue <- message
		log.Printf("Event in topic %s published - Message %v\n", topic, message)
	}()

	return nil
}

func (localP *localPubSub) Subscribe(ctx context.Context, topic pubsub.Topic) (ch <-chan *pubsub.Message, close func()) {
	localP.locker.Lock()
	chanSubMsg := make(chan *pubsub.Message)
	if listChannel, ok := localP.mapTopicChannel[topic]; ok {
		listChannel = append(listChannel, chanSubMsg)
	} else {
		localP.mapTopicChannel[topic] = []chan *pubsub.Message{chanSubMsg}
	}
	localP.locker.Unlock()

	return chanSubMsg, func() {
		log.Println("Unsubscribe")
		if listChan, ok := localP.mapTopicChannel[topic]; ok {
			for i, c := range listChan {
				if c == chanSubMsg {
					localP.locker.Lock()
					localP.mapTopicChannel[topic] = append(listChan[:i], listChan[i+1:]...)
					localP.locker.Unlock()
					break
				}
			}
		}
	}
}

func (localP *localPubSub) run() error {
	log.Println("Pubsub started")

	go func() {
		for {
			mess := <-localP.messageQueue
			log.Println("Message dequeued")

			if listChan, ok := localP.mapTopicChannel[mess.GetTopic()]; ok {
				for _, c := range listChan {
					go func(ch chan *pubsub.Message) {
						ch <- mess
					}(c)
				}
			}

		}
	}()

	return nil
}

func NewPubSubLocal() pubsub.PubSub {
	localP := &localPubSub{
		messageQueue:    make(chan *pubsub.Message, 100),
		mapTopicChannel: make(map[pubsub.Topic][]chan *pubsub.Message),
		locker:          new(sync.RWMutex),
	}

	_ = localP.run()
	return localP
}
