package datasetexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		"dataset",
		createDefaultConfig,
		exporter.WithLogs(createLogsExporter, component.StabilityLevelDevelopment),
		// TODO Should trace and metric exporters be added?
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		DatasetUrl: "https://app.scalyr.com",
	}
}

func createLogsExporter(ctx context.Context, set exporter.CreateSettings, config component.Config) (exporter.Logs, error) {
	cfg := config.(*Config)
	e, err := newDatasetExporter(cfg.ApiKey, cfg.DatasetUrl)
	if err != nil {
		return nil, err
	}

	return exporterhelper.NewLogsExporter(
		ctx,
		set,
		cfg,
		e.consumeLogs,
		// TODO What non-default options (retry, etc) should be used?
		//      The default is no queue, no retries and a timeout of five seconds
		//      Ref: https://pkg.go.dev/go.opentelemetry.io/collector/exporter/exporterhelper#Option
	)
}
