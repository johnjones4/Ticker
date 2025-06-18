package noaa

import (
	"context"
	"fmt"
	"log/slog"
	"main/ticker/core"
)

type NOAA struct {
	stationId string
	client    ClientWithResponsesInterface
	log       *slog.Logger
}

func (p *NOAA) Init(ctx context.Context, log *slog.Logger, cfg *core.Configuration) error {
	p.stationId = cfg.NoaaStationId
	p.log = log
	var err error
	p.client, err = NewClientWithResponses("https://api.weather.gov")
	if err != nil {
		return err
	}
	return nil
}

func (p *NOAA) Name() string {
	return "Natnl Weather Svc"
}

func (p *NOAA) Update(ctx context.Context) ([]string, error) {
	p.log.Info("updating NOAA")
	res, err := p.client.StationObservationLatestWithResponse(ctx, p.stationId, &StationObservationLatestParams{})
	if err != nil {
		return nil, err
	}
	if res.ApplicationgeoJSON200 == nil {
		return nil, nil
	}

	out := make([]string, 0)

	if res.ApplicationgeoJSON200.Properties.Temperature != nil && res.ApplicationgeoJSON200.Properties.Temperature.Value != nil {
		out = append(out, fmt.Sprintf("Temp: %0.1f°", celsiusToFahrenheit(*res.ApplicationgeoJSON200.Properties.Temperature.Value)))
	}

	if res.ApplicationgeoJSON200.Properties.WindSpeed != nil && res.ApplicationgeoJSON200.Properties.WindSpeed.Value != nil {
		out = append(out, fmt.Sprintf("Wind: %0.1fmph", kmToMiles(*res.ApplicationgeoJSON200.Properties.WindSpeed.Value)))
	}

	if res.ApplicationgeoJSON200.Properties.Dewpoint != nil && res.ApplicationgeoJSON200.Properties.Dewpoint.Value != nil {
		out = append(out, fmt.Sprintf("Dewpoint: %0.1f°", celsiusToFahrenheit(*res.ApplicationgeoJSON200.Properties.Dewpoint.Value)))
	}

	if res.ApplicationgeoJSON200.Properties.BarometricPressure != nil && res.ApplicationgeoJSON200.Properties.BarometricPressure.Value != nil {
		out = append(out, fmt.Sprintf("Pressure: %0.1finHg", millibarsToInHg(*res.ApplicationgeoJSON200.Properties.BarometricPressure.Value)))
	}

	if res.ApplicationgeoJSON200.Properties.RelativeHumidity != nil && res.ApplicationgeoJSON200.Properties.RelativeHumidity.Value != nil {
		out = append(out, fmt.Sprintf("Humidity: %0.0f%%", *res.ApplicationgeoJSON200.Properties.RelativeHumidity.Value))
	}

	return out, nil
}

func celsiusToFahrenheit(celsius float32) float32 {
	return (celsius * 9 / 5) + 32
}

func kmToMiles(kmh float32) float32 {
	return kmh * 0.621371
}

func millibarsToInHg(mb float32) float32 {
	return mb * 0.02953
}
