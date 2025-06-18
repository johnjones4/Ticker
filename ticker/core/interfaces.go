package core

import (
	"context"
	"log/slog"
)

type Init interface {
	Init(ctx context.Context, log *slog.Logger, cfg *Configuration) error
}

type Provider interface {
	Init
	Name() string
	Update(ctx context.Context) ([]string, error)
}

type Output interface {
	Init
	PrepareSegments(ctx context.Context, segs []string) error
	Update(ctx context.Context, msgs map[string][]string) error
}
