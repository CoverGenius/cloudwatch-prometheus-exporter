# cloudwatch-prometheus-exporter
Multi regional cloudwatch prometheus exporter

./cloudwatch-prometheus-exporter -config </path/to/your/config.yaml>

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
```
