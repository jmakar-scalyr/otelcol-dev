package datasetexporter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/collector/pdata/plog"
	"golang.org/x/time/rate"
)

type datasetExporter struct {
	apiKey     string
	datasetUrl string

	client  *http.Client
	limiter *rate.Limiter

	session string
}

func newDatasetExporter(apiKey, datasetUrl string) (*datasetExporter, error) {
	return &datasetExporter{
		apiKey:     apiKey,
		datasetUrl: datasetUrl,

		client:  &http.Client{Timeout: time.Second * 60},
		limiter: rate.NewLimiter(100*rate.Every(1*time.Minute), 100), // 100 requests / minute

		// TODO Use goroutines to support multiple sessions
		session: uuid.New().String(),
	}, nil
}

func (e *datasetExporter) consumeLogs(ctx context.Context, ld plog.Logs) error {
	// TODO Relying on default NewLogsExporter settings for queue, retries, timeouts
	//      The default is no queue, no retries and a timeout of five seconds
	//      Ref: https://pkg.go.dev/go.opentelemetry.io/collector/exporter/exporterhelper#Option

	// FIXME STOPPED Handle errors via retries, etc (using NewLogsExporter settings?)
	//               Explicitly log when batches are lost

	// FIXME Stop early if ctx if triggered
	// FIXME Record stats and output periodically (as events or log entries?)

	const MAX_BUFFER_SIZE = 6 * 1048576
	bufferPrefix := `{"session":"` + e.session + `","events":[`
	bufferSuffix := "]}"

	newBuffer := func() *bytes.Buffer {
		return bytes.NewBufferString(bufferPrefix)
	}

	marshalLog := func(log plog.LogRecord) ([]byte, error) {
		event := log.Attributes().AsRaw()

		if body := log.Body().Str(); body != "" {
			event["message"] = body
		}
		if dropped := log.DroppedAttributesCount(); dropped > 0 {
			event["droppedAttributesCount"] = dropped
		}
		if observed := log.ObservedTimestamp().AsTime(); !observed.Equal(time.Unix(0, 0)) {
			event["observedTimestamp"] = observed.String()
		}
		if sevNum := log.SeverityNumber(); sevNum > 0 {
			event["severityNumber"] = sevNum
		}
		if sevText := log.SeverityText(); sevText != "" {
			event["severityText"] = sevText
		}
		if span := log.SpanID().String(); span != "" {
			event["spanId"] = span
		}
		if timestamp := log.Timestamp().AsTime(); !timestamp.Equal(time.Unix(0, 0)) {
			event["timestamp"] = timestamp.String()
		}
		if trace := log.TraceID().String(); trace != "" {
			event["traceId"] = trace
		}

		if buf, err := json.Marshal(event); err != nil {
			return nil, err
		} else {
			// Timestamp is required, otherwise it appears be from zero Unix time (ie 1970)
			prefix := []byte(fmt.Sprintf(`{"ts":"%d000000000","attrs":`, time.Now().Unix()))
			suffix := []byte("}")
			return append(append(prefix, buf...), suffix...), nil
		}
	}

	sendBuffer := func(buf *bytes.Buffer) error {
		buf.WriteString(bufferSuffix)

		request, err := e.newRequest(ctx, "POST", e.datasetUrl+"/api/addEvents", bytes.NewBuffer(buf.Bytes()))
		if err != nil {
			return err
		}

		resp, err := e.client.Do(request)
		if err != nil {
			return err
		}
		resp.Body.Close()
		if !(200 <= resp.StatusCode && resp.StatusCode < 300) {
			return fmt.Errorf("unsuccessful (%d) status code", resp.StatusCode)
		}

		return nil
	}

	batchBuffer := newBuffer()

	resourceLogs := ld.ResourceLogs()
	for i := 0; i < resourceLogs.Len(); i++ {
		scopeLogs := resourceLogs.At(i).ScopeLogs()
		for j := 0; j < scopeLogs.Len(); j++ {
			logRecords := scopeLogs.At(j).LogRecords()
			for k := 0; k < logRecords.Len(); k++ {
				logRecord := logRecords.At(k)

				logBuffer, err := marshalLog(logRecord)
				if err != nil {
					return err
				}
				if len(logBuffer) > MAX_BUFFER_SIZE-(len(bufferPrefix)+len(bufferSuffix)) {
					return fmt.Errorf("event too long (%d bytes)", len(logBuffer))
				}

				if len(logBuffer)+batchBuffer.Len() > MAX_BUFFER_SIZE-len(bufferSuffix) {
					if err := sendBuffer(batchBuffer); err != nil {
						return err
					}
					batchBuffer = newBuffer()
				}
				batchBuffer.Write(logBuffer)
			}
		}
	}

	if batchBuffer.Len() > len(bufferPrefix) {
		if err := sendBuffer(batchBuffer); err != nil {
			return err
		}
	}

	return nil
}

func (e *datasetExporter) newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	if err := e.limiter.Wait(ctx); err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+e.apiKey)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "otel-datasetexporter/"+VERSION)

	return request, nil
}
