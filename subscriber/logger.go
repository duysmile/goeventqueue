package subscriber

import "log"

// Logger records errors that occur within jobs or subscribers.
type Logger interface {
	Error(msg string, err error)
}

type logger struct {
}

func (l logger) Error(msg string, err error) {
	log.Println(msg, err)
}

// NewDefaultLogger returns a basic Logger that writes errors to the standard
// log package.
func NewDefaultLogger() Logger {
	return &logger{}
}
