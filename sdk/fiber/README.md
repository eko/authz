# Authz GoFiber middleware

A [GoFiber](https://gofiber.io) (v2) middleware that authorizes requests against Authz.

## Installation

```bash
go get github.com/eko/authz/sdk/fiber
```

## Usage

```go
import (
    "github.com/eko/authz/sdk"
    authzfiber "github.com/eko/authz/sdk/fiber"
    "github.com/gofiber/fiber/v2"
)

authzClient, err := sdk.NewClient(&sdk.Config{
    ClientID:     "your-client-id",
    ClientSecret: "your-client-secret",
    GrpcAddr:     "localhost:8081",
})

app := fiber.New()

app.Put("/posts/:id", authzfiber.New(authzfiber.Config{
    Checker: authzClient,
    Principal: func(c *fiber.Ctx) (string, error) {
        // Return the principal identifier of the authenticated user,
        // for instance from a JWT or session set by a previous middleware.
        return c.Locals("user_id").(string), nil
    },
    ResourceKind:  "post",
    ResourceValue: authzfiber.ResourceValueFromParam("id"),
    Action:        "edit",
}), editPostHandler)
```

When `ResourceValue` is not set, the check is performed on the wildcard value `*` (useful for list endpoints).

By default, the middleware responds:

* `401 Unauthorized` when the principal cannot be identified,
* `403 Forbidden` when the principal is not allowed,
* `500 Internal Server Error` when the check cannot be performed (it fails closed).

Use the `ErrorHandler` option to customize these responses.
