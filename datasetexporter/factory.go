package datasetexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporterhelper"
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
		DataSetUrl: "https://app.scalyr.com/",
	}
}

func createLogsExporter(ctx context.Context, set exporter.CreateSettings, config component.Config) (exporter.Logs, error) {
	cfg := config.(*Config)

	// FIXME STOPPED
	return exporterhelper.NewLogsExporter(
		ctx,
		set,
		cfg,

	)
}
