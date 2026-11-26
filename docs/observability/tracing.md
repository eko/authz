# Observability: Tracing

Authz comes with observability instrumentation in order to have tracing enabled using [Jaeger](https://www.jaegertracing.io/) or any [OpenTelemetry Protocol (OTLP)](https://opentelemetry.io/docs/reference/specification/protocol/otlp/) compatible backend.

## How to enable

In order to enable tracing in Authz backend, you have to enable it using `APP_TRACE_ENABLED=true` environment variable.

You can also specify the following environment variables in order to configure it:

| Property | Default value | Description |
| -------- | ------------- | ----------- |
| APP_TRACE_ENABLED | `false` | Enable tracing observability using OpenTelemetry |
| APP_TRACE_EXPORTER | `jaeger` | Exporter you want to use. Could be `jaeger` or `otlpgrpc` |
| APP_TRACE_JAEGER_ENDPOINT | `localhost:14250` | Jaeger endpoint to be used |
| APP_TRACE_OTLP_DIAL_TIMEOUT | `3s` | OTLP gRPC exporter dial timeout value |
| APP_TRACE_OTLP_ENDPOINT | `localhost:30080` | OTLP gRPC endpoint value |
| APP_TRACE_SAMPLE_RATIO | `1.0` | Sampling ratio value defines how many traces should be sent to your exporter |

Once you will launch the backend, your traces should be collected over the exporter you provided.

## Zipkin exporter removal

The `zipkin` exporter has been removed because the OpenTelemetry Go Zipkin exporter is deprecated. The `APP_TRACE_ZIPKIN_URL` option has been removed too, and Authz now refuses to start when `APP_TRACE_EXPORTER=zipkin` is set.

To keep sending traces to Zipkin, use the OTLP exporter (`APP_TRACE_EXPORTER=otlpgrpc`) and set `APP_TRACE_OTLP_ENDPOINT` to an OTLP gRPC endpoint: a Zipkin server with OTLP ingestion enabled, or an [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/) that forwards traces to Zipkin. Jaeger also ingests OTLP natively.

## Try it using containers!

Following `docker-compose` files can be used in order to try these tools:

* [docker-compose.jaeger.yaml](https://github.com/eko/authz/blob/master/docker-compose.jaeger.yaml)
