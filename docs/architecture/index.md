# Architecture

Authz is a backend server for managing authorizations using RBAC or ABAC.

## How it works

<picture>
  <source media="(prefers-color-scheme: dark)" srcset="howitworks.dark.png">
  <img alt="Text changing depending on mode. Light: 'So light!' Dark: 'So dark!'" src="howitworks.png">
</picture>

Authz backend runs in the middle of your services infrastructure. It only communicates with its database.

Database can be `MySQL`, `PostgreSQL` or `SQLite`.

Your applications drive resources managed into Authz.

Applications on which you want to use authorizations can communicate over [`gRPC`](https://grpc.io/) by using the internal provided server (by default on port `8081`).

Elsewhere, you can also use the HTTP API (by default on port `8080`) to manage your Authz resources.

## APIs Documentation

gRPC remote procedure call methods are [available here](https://github.com/eko/authz/blob/master/backend/api/proto/api.proto)

HTTP OpenAPI documentation is [available here](https://editor.swagger.io/?raw=https://raw.githubusercontent.com/eko/authz/master/backend/internal/http/docs/swagger.yaml)
