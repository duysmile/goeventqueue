package eventqueue

type EventName string

type Event interface {
	GetName() EventName
	GetData() interface{}
}
