package goeventqueue

import (
	"context"
	"log"
	"time"
)

type JobConfig struct {
	MaxBackOff int64
}

type job struct {
	handler Handler
	backOff int64
	config  JobConfig
}

func (j *job) Run(ctx context.Context, data interface{}) error {
	var err error
	for {
		if err = j.handler(ctx, data); err != nil {
			log.Println("error handle job", err)
			j.backOff += 1
			if j.backOff >= j.config.MaxBackOff {
				return err
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(2 << j.backOff * time.Second):
			}
		} else {
			return nil
		}
	}
}

func NewJob(handler Handler, cfg JobConfig) *job {
	return &job{
		handler: handler,
		backOff: -1,
		config:  cfg,
	}
}
