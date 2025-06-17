package ticker

import (
	"context"
	"main/ticker/core"
	"sync"
	"time"
)

type Runtime struct {
	ConfigurationPath string
	Providers         []core.Provider
	Output            core.Output

	messages      chan []string
	configuration *core.Configuration
}

func (r *Runtime) Init(ctx context.Context) error {
	r.messages = make(chan []string, 128)

	r.configuration = &core.Configuration{}
	err := r.configuration.Load(r.ConfigurationPath)
	if err != nil {
		return err
	}

	inits := make([]core.Init, 0)
	for _, p := range r.Providers {
		inits = append(inits, p)
	}
	inits = append(inits, r.Output)

	for _, i := range inits {
		err := i.Init(ctx, r.configuration)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Runtime) Start(ctx context.Context) error {
	var wg sync.WaitGroup

	cancelableCtx, cancel := context.WithCancelCause(ctx)

	wg.Add(1)
	go func() {
		var currentMessages []string = []string{"Loading ..."}
		currentMessageIndex := 0
		displayTick := time.Tick(time.Second * 30)
		display := func() {
			if len(currentMessages) == 0 {
				return
			}
			err := r.Output.Display(cancelableCtx, currentMessages[currentMessageIndex%len(currentMessages)])
			if err != nil {
				cancel(err)
				return
			}
			currentMessageIndex++
		}
		display()
		for {
			select {
			case <-cancelableCtx.Done():
				wg.Done()
				return
			case msgs := <-r.messages:
				currentMessages = msgs
			case <-displayTick:
				display()
			}
		}
	}()

	update := time.Tick(time.Second)
	allMessages := make([][]string, len(r.Providers))
	for {
		select {
		case <-cancelableCtx.Done():
			wg.Wait()
			return cancelableCtx.Err()
		case <-update:
			hasUpdates := false
			for i, provider := range r.Providers {
				messages, err := provider.Update(cancelableCtx)
				if err != nil {
					cancel(err)
					continue
				}
				if messages != nil {
					allMessages[i] = messages
					hasUpdates = true
				}
			}
			if hasUpdates {
				out := make([]string, 0, len(r.Providers))
				for _, messages := range allMessages {
					out = append(out, messages...)
				}
				r.messages <- out
			}
		}
	}
}
