package youtube

import (
	"context"
	"fmt"
	"log/slog"
	"main/ticker/core"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type YouTube struct {
	service   *youtube.Service
	channelId string
	log       *slog.Logger
}

func (p *YouTube) Name() string {
	return "YouTube"
}

func (p *YouTube) Init(ctx context.Context, log *slog.Logger, cfg *core.Configuration) error {
	p.log = log
	p.channelId = cfg.YoutubeChannelId
	config, err := google.JWTConfigFromJSON(cfg.GoogleConfig, youtube.YoutubeReadonlyScope)
	if err != nil {
		return err
	}
	client := config.Client(ctx)
	service, err := youtube.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	p.service = service
	return nil
}

func (p *YouTube) Update(ctx context.Context) ([]string, error) {
	p.log.Info("Updating YouTube stats")
	res, err := p.service.Channels.List([]string{"statistics"}).Id(p.channelId).Do()
	if err != nil {
		return nil, err
	}
	if len(res.Items) == 0 || res.Items[0] == nil || res.Items[0].Statistics == nil {
		return nil, nil
	}
	return []string{
		fmt.Sprintf("Views: %s", core.FormatUintWithCommas(res.Items[0].Statistics.ViewCount)),
		fmt.Sprintf("Subscribers: %s", core.FormatUintWithCommas(res.Items[0].Statistics.SubscriberCount)),
	}, nil
}
