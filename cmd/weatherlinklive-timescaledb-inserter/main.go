package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/alpineworks/ootel"
	"github.com/michaelpeterswa/weatherlinklive-timescaledb-inserter/internal/config"
	"github.com/michaelpeterswa/weatherlinklive-timescaledb-inserter/internal/logging"
	"github.com/michaelpeterswa/weatherlinklive-timescaledb-inserter/internal/timescale"
	"github.com/tannerryan/davisweather"
)

func main() {
	slogHandler := slog.NewJSONHandler(os.Stdout, nil)
	slog.SetDefault(slog.New(slogHandler))

	slog.Info("welcome to weatherlinklive-timescaledb-inserter!")

	c, err := config.NewConfig()
	if err != nil {
		slog.Error("could not create config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slogLevel, err := logging.LogLevelToSlogLevel(c.String(config.LogLevel))
	if err != nil {
		slog.Error("could not parse log level", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.SetLogLoggerLevel(slogLevel)

	ctx := context.Background()

	ootelClient := ootel.NewOotelClient(
		ootel.WithMetricConfig(
			ootel.NewMetricConfig(
				c.Bool(config.MetricsEnabled),
				c.Int(config.MetricsPort),
			),
		),
		ootel.WithTraceConfig(
			ootel.NewTraceConfig(
				c.Bool(config.TracingEnabled),
				c.Float64(config.TracingSampleRate),
				c.String(config.TracingService),
				c.String(config.TracingVersion),
			),
		),
	)

	shutdown, err := ootelClient.Init(ctx)
	if err != nil {
		panic(err)
	}

	defer func() {
		_ = shutdown(ctx)
	}()

	timescaleClient, err := timescale.NewTimescaleClient(ctx, c.String(config.TimescaleConnString))
	if err != nil {
		slog.Error("could not create timescale client", slog.String("error", err.Error()))
		os.Exit(1)
	}

	wllClient, error := davisweather.Unmanaged(ctx, false, c.String(config.WeatherLinkLiveHost), 80)
	if error != nil {
		slog.Error("could not create weatherlink live client", slog.String("error", error.Error()))
		os.Exit(1)
	}

	for {
		<-wllClient.Notify
		report, err := wllClient.Report()
		if err != nil {
			slog.Error("could not get report", slog.String("error", err.Error()))
			continue
		}

		err = timescaleClient.Insert(ctx, report)
		if err != nil {
			slog.Error("could not insert report", slog.String("error", err.Error()))
			continue
		}

		slog.Debug("inserted report", slog.Any("device_id", report.DeviceID))
	}
}
