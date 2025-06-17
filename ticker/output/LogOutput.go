package output

import (
	"context"
	"log"
	"main/ticker/core"
)

type LogOutput struct {
}

func (o *LogOutput) Init(ctx context.Context, cfg *core.Configuration) error {
	return nil
}

func (o *LogOutput) Display(ctx context.Context, msg string) error {
	log.Printf("Message: %s", msg)
	return nil
}
