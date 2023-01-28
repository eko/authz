# How it works

Authz is a backend server for managing authorizations using RBAC or ABAC.

## Applications

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/eko/authz/master/docs/architecture/howitworks.dark.png">
  <img alt="Text changing depending on mode. Light: 'So light!' Dark: 'So dark!'" src="[howitworks.png](https://raw.githubusercontent.com/eko/authz/master/docs/architecture/howitworks.png)">
</picture>

Authz backend runs in the middle of your services infrastructure. It only communicates with its database.

Your applications drive resources managed into Authz.

Applications on which you want to use authorizations can communicate over [`gRPC`](https://grpc.io/) by using the internal provided server (by default on port `8081`).

Elsewhere, you can also use the HTTP API (by default on port `8080`) to manage your Authz resources.

## Storage / SQL Database

Actually, Authz supports the following SQL storages:

* [MySQL](https://www.mysql.com/)
* [PostgreSQL](https://www.postgresql.org)
* [SQLite](https://www.sqlite.org)

When running the backend, you can change SQL database driver using `DATABASE_DRIVER` environment variable.

## HTTP and gRPC APIs

We have documentations for our APIs: gRPC API is using [`Protocol Buffers`](https://developers.google.com/protocol-buffers?hl=fr) schema format and our HTTP API is using [OpenAPI](https://swagger.io/specification/) specification format.

* gRPC remote procedure call methods are [available here](https://github.com/eko/authz/blob/master/backend/api/proto/api.proto)
* HTTP OpenAPI documentation is [available here](https://editor.swagger.io/?url=https://raw.githubusercontent.com/eko/authz/master/backend/internal/http/docs/swagger.yaml)
