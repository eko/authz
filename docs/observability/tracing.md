# Observability: Tracing

Authz comes with observability instrumentation in order to have tracing enabled using [Jaeger](https://www.jaegertracing.io/), [Zipkin](https://zipkin.io/) or [OpenTelemetry Protocol (OTLP)](https://opentelemetry.io/docs/reference/specification/protocol/otlp/).

## How to enable

In order to enable tracing in Authz backend, you have to enable it using `APP_TRACE_ENABLED=true` environment variable.

You can also specify the following environment variables in order to configure it:

| Property | Default value | Description |
| -------- | ------------- | ----------- |
| APP_TRACE_ENABLED | `false` | Enable tracing observability using OpenTelemetry |
| APP_TRACE_EXPORTER | `jaeger` | Exporter you want to use. Could be `jaeger`, `zipkin` or `otlpgrpc` |
| APP_TRACE_JAEGER_URL | `http://localhost:14268/api/traces` | Jaeger API URL to be used |
| APP_TRACE_OTLP_DIAL_TIMEOUT | `3s` | OTLP gRPC exporter dial timeout value |
| APP_TRACE_OTLP_ENDPOINT | `localhost:30080` | OTLP gRPC endpoint value |
| APP_TRACE_SAMPLE_RATIO | `1.0` | Sampling ratio value defines how many traces should be sent to your exporter |
| APP_TRACE_ZIPKIN_URL | `http://localhost:9411/api/v2/spans` | Zipkin API URL to be used |

Once you will launch the backend, your traces should be collected over the exporter you provided.

## Try it using containers!

Following `docker-compose` files can be used in order to try these tools:

* [docker-compose.jaeger.yaml](https://github.com/eko/authz/blob/master/docker-compose.jaeger.yaml)
* [docker-compose.zipkin.yaml](https://github.com/eko/authz/blob/master/docker-compose.zipkin.yaml)