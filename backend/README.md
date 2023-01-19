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
| APP_AUDIT_CLEAN_DAYS_TO_KEEP | `7` | Audit logs number of days to keep in database |
| APP_AUDIT_CLEAN_DELAY | `1h` | Audit logs clean delay |
| APP_AUDIT_FLUSH_DELAY | `3s` | Delay in which audit logs will be batch into database |
| APP_AUDIT_RESOURCE_KIND_REGEX | `.*` | Filter which resource kind will be added on audit logs |
| APP_STATS_CLEAN_DAYS_TO_KEEP | `30` | Statistics number of days to keep in database |
| APP_STATS_CLEAN_DELAY | `1h` | Statistics clean delay |
| APP_STATS_FLUSH_DELAY | `3s` | Delay in which statistics will be batch into database |
| APP_STATS_RESOURCE_KIND_REGEX | `.*` | Filter which resource kind will be added on statistics |
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
| DISPATCHER_EVENT_CHANNEL_SIZE | `10000` | Event dispatcher channel size |
| GRPC_SERVER_ADDR | `:8081` | gRPC server address (hostname and port) |
| HTTP_SERVER_ADDR | `:8080` | HTTP server address (hostname and port) |
| HTTP_SERVER_CORS_ALLOW_CREDENTIALS | `true` | Should CORS allow credentials requests? |
| HTTP_SERVER_CORS_ALLOWED_DOMAINS | `http://localhost:3000` | CORS allowed domains |
| HTTP_SERVER_CORS_ALLOWED_HEADERS | `Authorization,Origin,Content-Length,Content-Type` | CORS allowed headers |
| HTTP_SERVER_CORS_ALLOWED_METHODS | `GET,POST,PATCH,PUT,DELETE,HEAD,OPTIONS` | CORS allowed methods |
| HTTP_SERVER_CORS_CACHE_MAX_AGE | `12h` | CORS cache max age value to be returned by server |
| LOGGER_LEVEL | `INFO` | Log level, could be `DEBUG`, `INFO`, `WARN` or `ERROR` |
| USER_ADMIN_DEFAULT_PASSWORD | `changeme` | Default admin password updated on app launch |
| OAUTH_CLIENT_ID | N/A | OAuth client ID provided by your issuer |
| OAUTH_CLIENT_SECRET | N/A | OAuth client Secret provider by your issuer |
| OAUTH_COOKIES_DOMAIN_NAME | `localhost` | OAuth domain name on which cookies will be stored |
| OAUTH_FRONTEND_REDIRECT_URL | `http://localhost:3000` | Frontend redirect URL when OAuth authentication is successful |
| OAUTH_ISSUER_URL | N/A | Issuer OpenID Connect URL (will be used to retrieve /.well-known/openid-configuration) |
| OAUTH_REDIRECT_URL | `[12h](http://localhost:8080/v1/oauth/callback)` | Backend OAuth callback URL |
| OAUTH_SCOPES | `profile,email` | OAuth scopes to be retrieved from your issuer |

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
