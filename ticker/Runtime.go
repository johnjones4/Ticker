package ticker

import (
	"context"
	"log/slog"
	"main/ticker/core"
	"main/ticker/output"
	"sync"
	"time"
)

type Runtime struct {
	ConfigurationPath string
	Providers         []core.Provider
	Log               *slog.Logger

	messages      chan map[string][]string
	configuration *core.Configuration
	output        core.Output
}

func (r *Runtime) Init(ctx context.Context) error {
	r.messages = make(chan map[string][]string, 128)

	r.configuration = &core.Configuration{}
	err := r.configuration.Load(r.ConfigurationPath)
	if err != nil {
		return err
	}

	if r.configuration.SerialDevice != nil {
		r.output = &output.LedSign{}
	} else {
		r.output = &output.LogOutput{}
	}

	segments := make([]string, 0, len(r.Providers))
	inits := make([]core.Init, 0)
	for _, p := range r.Providers {
		inits = append(inits, p)
		segments = append(segments, p.Name())
	}
	inits = append(inits, r.output)

	for _, i := range inits {
		err := i.Init(ctx, r.Log, r.configuration)
		if err != nil {
			return err
		}
	}

	err = r.output.PrepareSegments(ctx, segments)
	if err != nil {
		return err
	}

	return nil
}

func (r *Runtime) Start(ctx context.Context) error {
	var wg sync.WaitGroup

	cancelableCtx, cancel := context.WithCancelCause(ctx)

	wg.Add(1)
	go func() {
		for {
			select {
			case <-cancelableCtx.Done():
				wg.Done()
				return
			case msgs := <-r.messages:
				err := r.output.Update(cancelableCtx, msgs)
				if err != nil {
					cancel(err)
					continue
				}
			}
		}
	}()

	update := time.Tick(time.Second * 10)
	allMessages := make(map[string][]string)
	for {
		select {
		case <-cancelableCtx.Done():
			wg.Wait()
			return cancelableCtx.Err()
		case <-update:
			hasUpdates := false
			for _, provider := range r.Providers {
				messages, err := provider.Update(cancelableCtx)
				if err != nil {
					cancel(err)
					continue
				}
				if messages != nil {
					allMessages[provider.Name()] = messages
					hasUpdates = true
				}
			}
			if hasUpdates {
				r.Log.Info("Updated messages available")
				out := make(map[string][]string)
				for key, messages := range allMessages {
					out[key] = make([]string, len(messages))
					copy(out[key], messages)
				}
				r.messages <- out
			}
		}
	}
}
