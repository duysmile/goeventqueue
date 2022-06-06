package pubsub

import (
	"fmt"
	"time"
)

type Message struct {
	id        string
	topic     Topic
	data      interface{}
	createdAt time.Time
}

func NewMessage(data interface{}) *Message {
	now := time.Now()
	return &Message{
		id:        fmt.Sprintf("pubsub:%v", now.UnixNano()),
		data:      data,
		createdAt: now,
	}
}

func (m *Message) SetTopic(topic Topic) {
	m.topic = topic
}

func (m *Message) GetTopic() Topic {
	return m.topic
}

func (m *Message) GetData() interface{} {
	return m.data
}
