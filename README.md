# otelcol-dev
OpenTelemetry collector development environment

## Setup environment
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

## Add custom/development components
- Create and populate a directory with the new component, eg ./datasetexporter
  - Give the module an appropriate name and location
    - Eg `go mod init github.com/jmakar-scalyr/otelcol-dev/datasetexporter`
- In otelcol-dev, modify components.go to include the new component, eg:
  - ```sh
    $ diff -U2 components.go{.orig,}
    --- components.go.orig
    +++ components.go
    @@ -12,4 +12,5 @@
            batchprocessor "go.opentelemetry.io/collector/processor/batchprocessor"
            otlpreceiver "go.opentelemetry.io/collector/receiver/otlpreceiver"
    +       datasetexporter "github.com/jmakar-scalyr/otelcol-dev/datasetexporter"
     )

    @@ -33,4 +34,5 @@
            factories.Exporters, err = exporter.MakeFactoryMap(
                    loggingexporter.NewFactory(),
    +               datasetexporter.NewFactory(),
            )
            if err != nil {
    ```
- In otelcol-dev, modify go.mod to include the new requirement and associate it with a local path, eg:
  - ```sh
    $ diff -U1 go.mod{.orig,}
    --- go.mod.orig
    +++ go.mod
    @@ -6,2 +6,5 @@

    +require "github.com/jmakar-scalyr/otelcol-dev/datasetexporter" v0.0.0
    +replace "github.com/jmakar-scalyr/otelcol-dev/datasetexporter" v0.0.0 => "../datasetexporter"
    +
     require (
    ```
- In otelcol-dev, modify otel-config.yaml to include the new component, eg:
  - ```sh
    $ diff -U3 otel-config.yaml{.orig,}
    --- otel-config.yaml.orig
    +++ otel-config.yaml
    @@ -10,6 +10,9 @@
     exporters:
       logging:
         loglevel: debug
    +  dataset:
    +    apikey: <elided>
    +    dataseturl: https://app-qatesting.scalyr.com/

     service:
       pipelines:
    ```
- Build a new version with: `go get -u github.com/jmakar-scalyr/otelcol-dev/datasetexporter; go build -o otelcol-dev`
