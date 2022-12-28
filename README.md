# otelcol-dev
OpenTelemetry collector development environment

- Download ocb from https://github.com/open-telemetry/opentelemetry-collector/releases
  - `chmod +x ocb_...; xattr -d com.apple.quarantine ocb_...`
- Create builder-config.yaml
  - ```yaml
    dist:
      name: otelcol-dev
      output_path: ./otelcol-dev
     
    exporters:
      - import: go.opentelemetry.io/collector/exporter/loggingexporter
        gomod: go.opentelemetry.io/collector v0.64.0
     
    receivers:
      - import: go.opentelemetry.io/collector/receiver/otlpreceiver
        gomod: go.opentelemetry.io/collector v0.64.0
     
    processors:
      - import: go.opentelemetry.io/collector/processor/batchprocessor
        gomod: go.opentelemetry.io/collector v0.64.0
    ```
  - Ref: https://opentelemetry.io/docs/collector/custom-collector/#step-2---create-a-builder-manifest-file
- ./ocb_... --config builder-config.yaml
  - This creates otelcol-dev/ with source and binary
- Create otelcol-dev/otel-config.yaml
  - ```yaml
    receivers:
      otlp:
        protocols:
          grpc:
          http:
    
    processors:
      batch:
    
    exporters:
      logging:
        loglevel: debug
    
    service:
      pipelines:
        logs:
          receivers: [otlp]
          processors: [batch]
          exporters: [logging]
    ```
  - Ref: https://opentelemetry.io/docs/collector/configuration/
- Launch with: `cd otelcol-dev; otelcol-dev --config otel-config.yaml`
- Test with: `curl -i http://127.0.0.1:4318/v1/logs -H 'Content-Type: application/json' -d @test-log.json`
  - ```json
    {
      "resourceLogs": [
        {
          "scopeLogs": [
            {
              "logRecords": [
                {
                  "timeUnixNano": "1581452773000000789",
                  "body": {
                    "stringValue": "This is a log message"
                  },
                  "attributes": [
                    { 
                      "key": "app",
                      "value": {
                        "stringValue": "server"
                      }
                    },
                    {
                      "key": "instance_num",
                      "value": {
                        "intValue": "1"
                      }
                    }
                  ]
                }
              ]
            }
          ]
        }
      ]
    }
    ```
  - Ref: https://github.com/open-telemetry/opentelemetry-specification/blob/main/experimental/serialization/json.md
    - Note some docs have not been updated: https://github.com/open-telemetry/opentelemetry-collector/blob/main/CHANGELOG.md#-breaking-changes--18
