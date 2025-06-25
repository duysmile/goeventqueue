package goeventqueue

// EventName identifies the name of an event that will be published or
// consumed by the queue.
type EventName string

// Event represents a unit of information that can be passed through the
// queue. Implementations should provide the event name and its associated
// payload data.
type Event interface {
	GetName() EventName
	GetData() interface{}
}
