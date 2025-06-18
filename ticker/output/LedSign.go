package output

import (
	"context"
	"fmt"
	"log/slog"
	"main/ticker/core"
	"strings"
	"time"

	alphasign "github.com/johnjones4/alpha-sign-communications-protocol"
)

type LedSign struct {
	sign *alphasign.Sign
	log  *slog.Logger
}

var (
	textFileLabel alphasign.FileLabel = 'A'
)

func (o *LedSign) Init(ctx context.Context, log *slog.Logger, cfg *core.Configuration) error {
	o.log = log
	sign, err := alphasign.New(alphasign.SignAddressBroadcast, alphasign.AllSigns, *cfg.SerialDevice, 9600)
	if err != nil {
		return err
	}

	err = sign.Send(&alphasign.WriteSpecialFunctionCommand{
		Label: alphasign.ClearOrSetMemoryConfig,
		Data: []alphasign.Bytes{
			alphasign.MemoryConfiguration{
				FileLabel:                textFileLabel,
				FileType:                 alphasign.TextFile,
				KeyboardProtectionStatus: 'U',
				FileSize:                 alphasign.FileSize(1024),
			},
		},
	})
	if err != nil {
		return err
	}

	o.sign = sign
	return nil
}

func (o *LedSign) Update(ctx context.Context, msgs map[string][]string) error {
	strs := make([]string, 0, len(msgs))

	for label, msgs1 := range msgs {
		if len(msgs1) == 0 {
			continue
		}
		strs = append(strs, fmt.Sprintf("%s: %s", label, strings.Join(msgs1, ", ")))
	}

	strs = append(strs, fmt.Sprintf("Last Updated: %s", time.Now().Format(time.ANSIC)))

	msg := strings.Join(strs, " | ")
	o.log.Info("message", slog.Int("size", len([]byte(msg))), slog.String("message", msg))

	return o.sign.Send(alphasign.WriteTextCommand{
		FileLabel: textFileLabel,
		Mode: &alphasign.TextMode{
			DisplayPosition: alphasign.Left,
			ModeCode:        alphasign.Rotate,
		},
		Message: append([]byte{0x15, 0x1C, 0x31, 0x10}, []byte(msg)...),
	})
}
