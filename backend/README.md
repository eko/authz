Authz - Backend
===============

This is the backend server of Authz.

Written in Go, it brings:

* an HTTP API server (default port: 8080)
* a gRPC API server (default port: 8081)

## Pre-requisites

In order to launch the backend server, you need to have a database running. Please refer to root [`README.md`](https://github.com/eko/authz) file.

## How to run

You can simply run it with:

```bash
$ go run cmd/main.go
```

## Configuration

Here are the available configuration options available as environment variable:

| Property | Default value | Description |
| -------- | ------------- | ----------- |
| APP_STATS_FLUSH_DELAY | `5s` | Delay in which statistics will be batch into database |
| AUTH_ACCESS_TOKEN_DURATION | `6h` | Access token duration  |
| AUTH_DOMAIN | `http://localhost:8080` | OAuth domain to be used  |
| AUTH_JWT_SIGN_STRING | `4uthz-s3cr3t-valu3-pl3as3-ch4ng3!` | Default HMAC to use for JWT tokens |
| AUTH_REFRESH_TOKEN_DURATION | `6h` | Refresh token duration |
| DATABASE_DRIVER | `postgres` | Database driver (`mysql`, `postgres` or `sqlite`) |
| DATABASE_HOST | `localhost` | Database host |
| DATABASE_NAME | `root` | Database name |
| DATABASE_PASSWORD | `toor` | Database password |
| DATABASE_PORT | `5432` | Database port |
| DATABASE_SSL | `disable` | Should database SSL mode be enabled? |
| DATABASE_TIMEZONE | `UTC` | Database timezone for date/time |
| DATABASE_USER | `root` | Database user |
| GRPC_SERVER_ADDR | `:8081` | gRPC server address (hostname and port) |
| HTTP_SERVER_ADDR | `:8080` | HTTP server address (hostname and port) |
| HTTP_SERVER_CORS_ALLOW_CREDENTIALS | `true` | Should CORS allow credentials requests? |
| HTTP_SERVER_CORS_ALLOWED_DOMAINS | `http://localhost:3000` | CORS allowed domains |
| HTTP_SERVER_CORS_ALLOWED_HEADERS | `Authorization,Origin,Content-Length,Content-Type` | CORS allowed headers |
| HTTP_SERVER_CORS_ALLOWED_METHODS | `GET,POST,PATCH,PUT,DELETE,HEAD,OPTIONS` | CORS allowed methods |
| HTTP_SERVER_CORS_CACHE_MAX_AGE | `12h` | CORS cache max age value to be returned by server |
| LOGGER_LEVEL | `INFO` | Log level, could be `DEBUG`, `INFO`, `WARN` or `ERROR` |
| USER_ADMIN_DEFAULT_PASSWORD | `changeme` | Default admin password updated on app launch |

## Tests

### Unit tests

You can run Go unit tests with:

```bash
$ make test-unit
```

### Functional tests

You can run HTTP API functional tests with:

```bash
$ make test-functional [tags=<something>]
```

Giving a tag is optionnal, it allows to run a specific tagged resource only.
