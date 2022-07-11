package goeventqueue

import "log"

type Logger interface {
	Error(msg string, err error)
}

type logger struct {
}

func (l logger) Error(msg string, err error) {
	log.Println(msg, err)
}

func NewDefaultLogger() Logger {
	return &logger{}
}
