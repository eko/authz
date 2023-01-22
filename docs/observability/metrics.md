# Observability: Metrics

Authz comes with observability instrumentation and some metrics you can retrieve in your [Prometheus](https://prometheus.io/) timeseries database.

## How to enable

In order to enable metrics in Authz backend, you have to simple enable it using `APP_METRICS_ENABLED=true` environment variable.

Once you will launch the backend, your metrics should be able to be scrapped using `/v1/metrics` endpoint.

## Available metrics

The following metrics are available at this time:

| Metric name | Labels | Description |
| ----------- | ------ | ----------- |
| `authz_check_counter` | `is_allowed`, `resource_kind` | The total number of checks processed |
| `authz_item_counter` | `item_type`, `action` | The total number of items (resource, policy, ...) created or updated in database |