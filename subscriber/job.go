package subscriber

import (
	"context"
	"time"
)

// JobConfig controls the retry behaviour when executing a job handler.
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

// NewJob creates a job that will run the given handler applying the retry
// configuration provided.
func NewJob(handler Handler, cfg JobConfig) *job {
	return &job{
		handler: handler,
		backOff: -1,
		config:  cfg,
	}
}
