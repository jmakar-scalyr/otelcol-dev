package datasetexporter

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/pdata/plog"
)

type datasetExporter struct {
	apiKey     string
	datasetUrl string

	marshaler  *plog.JSONMarshaler

	// FIXME rate limiter
}

func newDatasetExporter(apiKey, datasetUrl string) (*datasetExporter, error) {
	return &datasetExporter{
		apiKey: apiKey,
		datasetUrl: datasetUrl,
		marshaler: &plog.JSONMarshaler{},
	}, nil
}

func (e *datasetExporter) consumeLogs(ctx context.Context, ld plog.Logs) error {
	buf, err := e.marshaler.MarshalLogs(ld)
	if err != nil {
		return err
	}

	// FIXME STOPPED
	fmt.Printf("%s\n", string(buf))
	return nil
}
