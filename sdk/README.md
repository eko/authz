# Authz Go SDK

This is the Authz development kit for Go.

## Installation

You can install in your projects by importing the following dependency:

```bash
$ go get github.com/eko/authz/sdk@latest
```

## Usage

You have to instanciate a new Authz Client in your code by doing:

```go
authzClient, err := sdk.NewClient(&sdk.Config{
    ClientID: "your-client-id",
    ClientSecret: "your-client-secret",
    GrpcAddr: "localhost:8081",
})
```

Once the client is instanciate, you have access to all the gRPC methods and also some overridden ones.

In order to create a new Principal, you can use

```go
response, err := authzClient.PrincipalCreate(ctx, &authz.PrincipalCreateRequest{
    Id: "user-123",
    Attributes: []*authz.Attribute{
        {Key: "email", Value: "johndoe@acme.tld"},
    },
})
```

To declare a new resource:

```go
response, err := authzClient.ResourceCreate(ctx, &authz.ResourceCreateRequest{
    Id: "post.456",
    Kind: "post",
    Value: "456",
    Attributes: []*authz.Attribute{
        {Key: "owner_email", Value: "johndoe@acme.tld"},
    },
})
```

You can also declare a new policy this way:

```go
import (
    "github.com/eko/authz/backend/sdk/rule"
)

response, err := authzClient.PolicyCreate(ctx, &authz.PolicyCreateRequest{
    Id: "post-owners",
    Resources: []string{"post.*"},
    Actions: []string{"edit", "delete"},
    AttributeRules: []string{
        rule.AttributeEqual(
            rule.PrincipalResourceAttribute{
                PrincipalAttribute: "email",
                ResourceAttribute:  "owner_email",
            },
        ),
    },
})
```

Then, you can perform a check with:

```go
isAllowed, err := authzClient.IsAllowed(&authz.Check{
    Principal: "user-123",
    ResourceKind: "post",
    ResourceValue: "456",
    Action: "edit",
})
if err != nil {
    // Log error
}

if isAllowed {
    // Do something
}
```

Please note that you have access to all the gRPC methods [declared here](https://github.com/eko/authz/blob/master/backend/api/proto/api.proto) in the proto file.

## Configuration

This SDK connects over gRPC to the backend service. Here are the available configuration options:

| Property | Default value | Description |
| -------- | ------------- | ----------- |
| ClientID | *None* | Your service account client id used to authenticate |
| ClientSecret | *None* | Your service account client secret key used to authenticate |
| GrpcAddr | 127.0.0.1:8081 | Authz backend to connect to |

## Test

Unit tests can be run with:

```go
$ go test -v -race -count=1 ./...
```