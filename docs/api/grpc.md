# gRPC API

You can use the HTTP API on default port `8081`.

Protocol Buffers schema is [available here](https://github.com/eko/authz/blob/master/backend/api/proto/api.proto).

## Authentication

In order to authenticate with gRPC API, you will need to create a new service account.

You can create it on frontend or using the HTTP API `/v1/clients` endpoint.

Then, you can perform an authentication this way (here we are using `grpcurl` CLI tool):

```bash
grpcurl \
  -proto api.proto \
  -import-path $GOPATH/src/github.com/eko/authz/backend/api/proto/ \
  -plaintext \
  -d '{"client_id": "3d251d2c-9fe8-11ed-9cde-acde48001122", "client_secret": "JQ9IWwQBot9kNATVhacDpDhDqOBW7HZbIqAKh__LdwEP4xvH"}' \
  127.0.0.1:8081 \
  authz.Api/Authenticate

{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhdXRoeiIsInN1YiI6ImF1dGh6LXNhLXRlc3Qtc2EiLCJleHAiOjE2NzUwMjk1NjQsImlhdCI6MTY3NTAwNzk2NH0.Y0PrIdOMlE4Id07HYaLGDc9wDVG56OulX7iwNhdLPxs",
  "type": "bearer",
  "expiresIn": "21599"
}
```

You will gain an access token, available for 1 hour, by default.

Then, you can call the available RPC methods by adding a `Authorization: Bearer <token>` HTTP header to the request.

For instance:

```bash
grpcurl \
  -proto api.proto \
  -import-path $GOPATH/src/github.com/eko/authz/backend/api/proto/ \
  -plaintext \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"id": "user-456", "attributes": [{"key": "test1", "value": "value1"}]}' \
  127.0.0.1:8081 \
  authz.Api/PrincipalCreate

{
  "principal": {
    "id": "user-456",
    "attributes": [
      {
        "key": "test1",
        "value": "value1"
      }
    ]
  }
}
```