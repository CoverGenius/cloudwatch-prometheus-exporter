# cloudwatch-prometheus-exporter

[![Go Report Card](https://goreportcard.com/badge/github.com/CoverGenius/cloudwatch-prometheus-exporter)][goreportcard]

Cloudwatch prometheus exporter with support for multiple regions and counter metrics

`./cloudwatch-prometheus-exporter -config </path/to/your/config.yaml>`

## Configuration

The exporter accepts a YAML configuration file, the location of which can be specified by the `-config` flag.

Example of config file:
```
---
regions:
  - eu-central-1
  - ap-southeast-2
api_key: <API_KEY>
api_secret: <API_SECRET>
tags:
  - name: Environment
    value: production
period: 5
poll_interval: 5
log_level: 4
metrics:
  AWS/S3:
    length: 10080
```

Name          | Description
--------------|------------
listen        | Optional. Address for the Prometheus HTTP API to listen on. Defaults to `127.0.0.1:8080`
regions       | Required. List of AWS regions to query resources/metrics for.
api_key       | Required. AWS API Key ID.
api_secret    | Required. AWS API Secret.
tags          | Optional. List of name, value pairs used to filter AWS resources.
period        | Optional. How far back to request data for in minutes. Defaults to 5 minutes
poll_interval | Optional. How often in minutes to fetch new data from the Cloudwatch API, should be less than or equal to Period. Defaults to 5 minutes.
log_level     | Optional. Logging verbosity, must be between 1 and 5 inclusive. Higher levels represent greater verbosity. Defaults to 3 (log warnings and above).
metrics       | Optional. Map of configuration overrides for metric namespaces, currently only supports overriding `length` (equivalent to `period`).

[goreportcard]: https://goreportcard.com/report/github.com/CoverGenius/cloudwatch-prometheus-exporter
