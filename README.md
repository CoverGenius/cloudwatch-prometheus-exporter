# Cloudwatch Prometheus Exporter

[![Go Report Card](https://goreportcard.com/badge/github.com/CoverGenius/cloudwatch-prometheus-exporter)][goreportcard]

CloudWatch prometheus exporter with support for multiple regions and counter metrics.

`./cloudwatch-prometheus-exporter -config </path/to/your/config.yaml>`

## Configuration

The exporter accepts a YAML configuration file, the location of which can be specified by the `-config` flag.

Example configuration file:
```yaml
---
regions:
  - eu-central-1
  - ap-southeast-2
api_key: <API_KEY>
api_secret: <API_SECRET>
tags:
  - name: Environment
    value: production
period_seconds: 60
range_seconds: 300
poll_interval: 60
log_level: 4
metrics:
  AWS/EC2:
  AWS/RDS:
      - metric: CPUCreditBalance
      - metric: WriteIOPS
  AWS/S3:
      - metric: BucketSizeBytes
        period_seconds: 86400
        range_seconds: 604800
        dimensions:
            - name: StorageType
              value: StandardStorage
        statistics: [Maximum]
```

### Top level options

Name              | Description
------------------|------------
`listen`          | Optional. Address for the Prometheus HTTP API to listen on. Defaults to `127.0.0.1:8080`
`regions`         | Required. List of AWS regions to query resources/metrics for.
`api_key`         | Required. AWS API Key ID.
`api_secret`      | Required. AWS API Secret.
`tags`            | Optional. List of name, value pairs used to filter AWS resources.
`poll_interval`   | Optional. How often in seconds to fetch new data from the CloudWatch API, should be less than or equal to Period. Defaults to 300 (5 minutes).
`log_level`       | Optional. Logging verbosity, must be between 1 and 5 inclusive. Higher levels represent greater verbosity. Defaults to 3 (log warnings and above).
`period_seconds`  | Optional. Granularity of data retrieved from CloudWatch. Defaults to 60 (1 minute).
`range_seconds`   | Optional. How far back to request data for in seconds. Defaults to 300 (5 minutes).
`metrics`         | Optional. Map of metric configurations keyed by CloudWatch namespace, see per metric options below.

### Per metric options

The exporter will not query metrics for a namespace unless there is a key for that namespace under the `metrics` option. If only the namespace key is set then the default metrics for that namespace will be used. Otherwise individual metrics can be configured using the options below.

Name              | Description
------------------|------------
`metric`          | Required. Cloudwatch metric to use.
`output_name`     | Optional. Name to use for the generated Prometheus metric. Defaults to `<snake_case_metric>_<statistic>` if not set.
`help`            | Optional. The help text to use for the generated Prometheus metric. Defaults are configured for most CloudWatch metrics.
`statistics`      | Optional. List of CloudWatch statistics to generate metric series for. Defaults to `[Average]`.
`period_seconds`  | Optional. Granularity of data retrieved from CloudWatch. Defaults to 60 (1 minute).
`range_seconds`   | Optional. How far back to request data for in seconds. Defaults to global `range_seconds` if not set.

[goreportcard]: https://goreportcard.com/report/github.com/CoverGenius/cloudwatch-prometheus-exporter
