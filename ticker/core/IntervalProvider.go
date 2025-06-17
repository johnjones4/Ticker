package core

import (
	"context"
	"time"
)

type IntervalProvider struct {
	Parent         Provider
	UpdateInterval time.Duration
	lastUpdate     time.Time
}

func (ip *IntervalProvider) Init(ctx context.Context, cfg *Configuration) error {
	return ip.Parent.Init(ctx, cfg)
}

func (ip *IntervalProvider) Update(ctx context.Context) ([]string, error) {
	nextUpdate := ip.lastUpdate.Add(ip.UpdateInterval)
	if nextUpdate.After(time.Now()) {
		return nil, nil
	}
	msgs, err := ip.Parent.Update(ctx)
	if err != nil {
		return nil, err
	}
	ip.lastUpdate = time.Now()
	return msgs, nil
}
