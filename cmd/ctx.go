package cmd

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
)

func buildContext(ctx context.Context) (context.Context, context.CancelFunc) {
	if Timeout == "" {
		return context.WithCancel(ctx)
	}

	timeoutFlag, err := time.ParseDuration(Timeout)
	if err != nil {
		logrus.Warnf("failed to parse timeout value %v", err)
		return context.WithCancel(ctx)
	}

	return context.WithTimeout(ctx, timeoutFlag)
}
