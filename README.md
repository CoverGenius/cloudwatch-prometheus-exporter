# cloudwatch-prometheus-exporter

[![Go Report Card](https://goreportcard.com/badge/github.com/CoverGenius/cloudwatch-prometheus-exporter)][goreportcard]

Cloudwatch prometheus exporter with support for multiple regions and counter metrics

`./cloudwatch-prometheus-exporter -config </path/to/your/config.yaml>`

## Configuration

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

[goreportcard]: https://goreportcard.com/report/github.com/CoverGenius/cloudwatch-prometheus-exporter
