package output

import (
	"context"
	"log/slog"
	"main/ticker/core"
)

type LogOutput struct {
	log *slog.Logger
}

func (o *LogOutput) Init(ctx context.Context, log *slog.Logger, cfg *core.Configuration) error {
	o.log = log
	return nil
}

func (o *LogOutput) PrepareSegments(ctx context.Context, segs []string) error {
	return nil
}

func (o *LogOutput) Update(ctx context.Context, msgs map[string][]string) error {
	for i, msgs1 := range msgs {
		for _, msg := range msgs1 {
			o.log.Info("message", slog.String("n", i), slog.String("message", msg))
		}
	}
	return nil
}
