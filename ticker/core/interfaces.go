package core

import "context"

type Init interface {
	Init(ctx context.Context, cfg *Configuration) error
}

type Provider interface {
	Init
	Update(ctx context.Context) ([]string, error)
}

type Output interface {
	Init
	Display(ctx context.Context, msg string) error
}
