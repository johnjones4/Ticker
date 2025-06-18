package output

import (
	"context"
	"fmt"
	"log/slog"
	"main/ticker/core"
	"slices"
	"strings"
	"time"

	alphasign "github.com/johnjones4/alpha-sign-communications-protocol"
)

type LedSign struct {
	sign            *alphasign.Sign
	log             *slog.Logger
	orderedSegments []string
}

const (
	textFileLabel        alphasign.FileLabel = 0x31
	stringFileLabelStart alphasign.FileLabel = 0x32

	stringFileWidth = 125
	nStringFiles    = (2048 / stringFileWidth) + 1
)

func (o *LedSign) Init(ctx context.Context, log *slog.Logger, cfg *core.Configuration) error {
	o.log = log
	sign, err := alphasign.New(alphasign.SignAddressBroadcast, alphasign.AllSigns, *cfg.SerialDevice, 9600)
	if err != nil {
		return err
	}
	o.sign = sign

	configs := []alphasign.Bytes{
		alphasign.MemoryConfiguration{
			FileLabel:                textFileLabel,
			FileType:                 alphasign.TextFile,
			KeyboardProtectionStatus: 'U',
			FileSize:                 alphasign.FileSize(1024),
		},
	}
	displayString := []byte{0x15, 0x1C, 0x31}

	for i := range nStringFiles {
		filename := stringFileLabelStart + alphasign.FileLabel(i)
		o.log.Info("creating string file", slog.String("filename", fmt.Sprint(filename)))
		displayString = append(displayString, 0x10, byte(filename))
		configs = append(configs, alphasign.MemoryConfiguration{
			FileLabel:                filename,
			FileType:                 alphasign.StringFile,
			KeyboardProtectionStatus: 'U',
			FileSize:                 alphasign.FileSize(stringFileWidth),
		})
	}

	err = o.sign.Send(&alphasign.WriteSpecialFunctionCommand{
		Label: alphasign.ClearOrSetMemoryConfig,
		Data:  configs,
	})
	if err != nil {
		return err
	}

	err = o.sign.Send(alphasign.WriteTextCommand{
		FileLabel: textFileLabel,
		Mode: &alphasign.TextMode{
			DisplayPosition: alphasign.Left,
			ModeCode:        alphasign.Rotate,
		},
		Message: displayString,
	})

	if err != nil {
		return err
	}

	return nil
}

func (o *LedSign) PrepareSegments(ctx context.Context, segs []string) error {
	o.orderedSegments = segs
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

	n := 0
	for chunk := range slices.Chunk([]byte(msg), stringFileWidth) {
		err := o.sign.Send(alphasign.WriteStringCommand{
			FileLabel: stringFileLabelStart + alphasign.FileLabel(n),
			FileData:  chunk,
		})
		if err != nil {
			return err
		}
		n++
	}

	if n < nStringFiles {
		for i := range nStringFiles - n {
			err := o.sign.Send(alphasign.WriteStringCommand{
				FileLabel: stringFileLabelStart + alphasign.FileLabel(n+i),
				FileData:  []byte{},
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
