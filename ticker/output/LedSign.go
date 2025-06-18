package output

import (
	"context"
	"main/ticker/core"

	alphasign "github.com/johnjones4/alpha-sign-communications-protocol"
)

type LedSign struct {
	sign *alphasign.Sign
}

func (o *LedSign) Init(ctx context.Context, cfg *core.Configuration) error {
	sign, err := alphasign.New(alphasign.SignAddressBroadcast, alphasign.AllSigns, *cfg.SerialDevice, 9600)
	if err != nil {
		return err
	}
	o.sign = sign
	return nil
}

func (o *LedSign) Display(ctx context.Context, msg string) error {
	return o.sign.Send(alphasign.WriteTextCommand{
		FileLabel: 'A',
		Mode: &alphasign.TextMode{
			DisplayPosition: alphasign.Fill,
			ModeCode:        alphasign.Hold,
		},
		Message: []byte(msg),
	})
}
