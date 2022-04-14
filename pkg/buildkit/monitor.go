package buildkit

import (
	"context"

	"github.com/moby/buildkit/client"
	"github.com/sirupsen/logrus"
)

type Monitor interface {
	Monitor(ctx context.Context, ch chan *client.SolveStatus) error
}

type generalMonitor struct {
}

func NewMonitor() Monitor {
	return &generalMonitor{}
}

func (g generalMonitor) Monitor(ctx context.Context, ch chan *client.SolveStatus) error {
	for {
		select {
		case ss, ok := <-ch:
			if !ok {
				return nil
			}
			for _, vs := range ss.Statuses {
				logrus.Debug(vs)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
