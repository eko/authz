# HTTP API

You can use the HTTP API on default port `8080`.

OpenAPI specifications are available at [`http://localhost:8080/v1/swagger`](http://localhost:8080/v1/swagger).

You can also access [OpenAPI specifications](https://editor.swagger.io/?url=https://raw.githubusercontent.com/eko/authz/master/backend/internal/http/docs/swagger.yaml) here.

## Authentication

You can authenticate with the API using both username/password or a service account using OAuth `client_credentials` grant type.

### Username / Password

In order to authenticate with Username / Password, you can use this endpoint:

```bash
$ curl -X POST \
  -H 'Content-Type: application/json' \
  -d '{"username": "admin", "password": "changeme"}' \
  http://localhost:8080/v1/a

{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhdXRoeiIsInN1YiI6ImFkbWluIiwiZXhwIjoxNjc1MDI2ODM3LCJpYXQiOjE2NzUwMDUyMzd9.4oIMXD3jq6fe0Ykh3qcsBLlXG6vHLrOEo3xzD11LixM",
  "token_type": "bearer",
  "expires_in": 21599,
  "user": {
    "username": "admin",
    "created_at": "2023-01-29T16:12:20.671231+01:00",
    "updated_at": "2023-01-29T16:12:20.671231+01:00"
  }
}
```

Access token will be retrieved with the `expires_in` value indicating when it will expire. You will also have user information under the `user` key.

### Service account

You will first need to create a service account in order to use this authentication mode:

```bash
$ curl -X POST \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"name": "test-sa"}' \
  http://localhost:8080/v1/clients

{
  "client_id":"3d251d2c-9fe8-11ed-9cde-acde48001122",
  "client_secret":"JQ9IWwQBot9kNATVhacDpDhDqOBW7HZbIqAKh__LdwEP4xvH",
  "name":"test-sa",
  "domain":"http://localhost:8080",
  "created_at":"2023-01-29T16:18:52.542633+01:00",
  "updated_at":"2023-01-29T16:18:52.542633+01:00"
}
```

Once you have the `client_id` and `client_secret` informations, you can authenticate using the `/v1/token` endpoint:

```bash
$ curl -s -X POST -H 'Content-type: application/x-www-form-urlencoded' -d 'client_id=3d251d2c-9fe8-11ed-9cde-acde48001122&client_secret=JQ9IWwQBot9kNATVhacDpDhDqOBW7HZbIqAKh__LdwEP4xvH&grant_type=client_credentials' http://localhost:8080/v1/token | jq

{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJhdXRoeiIsInN1YiI6IjNkMjUxZDJjLTlmZTgtMTFlZC05Y2RlLWFjZGU0ODAwMTEyMiIsImV4cCI6MTY3NTAyNzMwNCwiaWF0IjoxNjc1MDA1NzA0fQ.a60XGmZBqVWjKc8t10ESp1tFrEH9XCFk5xtcwIZEM5Q",
  "expires_in": 21600,
  "refresh_token": "NJZMMWNKMWUTMGIXNY01NJA3LWE0ZTETM2VLYTHLNGRKM2UY",
  "token_type": "Bearer"
}
```

You will obtain an access token and also a refresh token in case you want to refresh the token and avoid it to expire.

## Check

You can perform a check using the following API (**you can perform multiple checks in a single request**):

```bash
 curl -X POST \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"checks": [{"principal": "user-123", "resource_kind": "post", "resource_value": "123", "action": "edit"}]}' \
  http://localhost:8080/v1/check | jq

{
  "checks": [
    {
      "principal": "user-123",
      "resource_kind": "post",
      "resource_value": "123",
      "action": "edit",
      "is_allowed": true
    }
  ]
}
```

## Policy

A policy allows to give access to one or multiple resources to perform one or multiple actions.

A policy can also contains optional `attribute_rules` if you want to use Attribute-Based Access Control.

### Declare a new policy

A policy can be created using the following API

```bash
$ curl -s -X POST \
  -H 'Content-type: application/json' \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"id": "post-manage", "resources": ["post.*"], "actions": ["edit", "delete"], "attribute_rules": ["resource.owner_email == principal.email"]}' \
  http://localhost:8080/v1/policies | jq

{
  "id": "post-manage",
  "resources": [
    {
      "id": "post.*",
      "kind": "post",
      "value": "*",
      "is_locked": false,
      "created_at": "2023-01-29T16:31:05.655777+01:00",
      "updated_at": "2023-01-29T16:31:05.655777+01:00"
    }
  ],
  "actions": [
    {
      "id": "edit",
      "created_at": "2023-01-29T16:31:05.656086+01:00",
      "updated_at": "2023-01-29T16:31:05.656086+01:00"
    },
    {
      "id": "delete",
      "created_at": "2023-01-29T16:12:20.665704+01:00",
      "updated_at": "2023-01-29T16:12:20.665704+01:00"
    }
  ],
  "attribute_rules": [
    "resource.owner_email == principal.email"
  ],
  "created_at": "2023-01-29T16:31:59.05117+01:00",
  "updated_at": "2023-01-29T16:31:59.05117+01:00"
}

```

### Delete a policy

```bash
$ curl -s -X DELETE \
-H "Authorization: Bearer $ACCESS_TOKEN" \
http://localhost:8080/v1/policies/post-manage | jq
{
  "success": true
}
```

## Principal

A principal is someone or something that can perform actions on resources. It can be a user, a group of user, an application or anything you want.

It just contains an identifier and optionnally a set of attributes.

### Declare a new principal

A principal can be created using the following API (attributes are optional)

```bash
$ curl -s -X POST \
  -H 'Content-type: application/json' \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"id": "user-123", "attributes": [{"key": "email", "value": "johndoe@acme.tld"}]}' \
  http://localhost:8080/v1/principals | jq

{
  "id": "user-123",
  "attributes": [
    {
      "key": "email",
      "value": "johndoe@acme.tld"
    }
  ],
  "is_locked": false,
  "created_at": "2023-01-29T16:25:56.947229+01:00",
  "updated_at": "2023-01-29T16:25:56.947229+01:00"
}
```

### Delete a principal

```bash
$ curl -s -X DELETE \
-H "Authorization: Bearer $ACCESS_TOKEN" \
http://localhost:8080/v1/principals/user-123 | jq
{
  "success": true
}
```

## Resource

A resource is every object you can do actions on. For instance, a blog post is a resource.

### Declare a new resource

A resource can be created using the following API (attributes are optional)

```bash
$ curl -s -X POST \
  -H 'Content-type: application/json' \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"id": "post.123", "kind": "post", "value": "123", "attributes": [{"key": "owner_email", "value": "johndoe@acme.tld"}]}' \
  http://localhost:8080/v1/resources | jq

{
  "id": "post.123",
  "kind": "post",
  "value": "123",
  "attributes": [
    {
      "key": "owner_email",
      "value": "johndoe@acme.tld"
    }
  ],
  "is_locked": false,
  "created_at": "2023-01-29T16:28:43.663332+01:00",
  "updated_at": "2023-01-29T16:28:43.663332+01:00"
}
```

### Delete a resource

```bash
$ curl -s -X DELETE \
-H "Authorization: Bearer $ACCESS_TOKEN" \
http://localhost:8080/v1/resources/post.123 | jq
{
  "success": true
}
```

## Role

If you want to use Role-Based Access Control, you can manage roles using these APIs.

### Declare a new role

```bash
$ curl -s -X POST \
  -H 'Content-type: application/json' \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"id": "post-manager", "policies": ["post-manage"]}' \
  http://localhost:8080/v1/roles | jq

{
  "id": "post-manager",
  "policies": [
    {
      "id": "post-manage",
      "attribute_rules": [
        "resource.owner_email == principal.email"
      ],
      "created_at": "2023-01-29T16:31:59.05117+01:00",
      "updated_at": "2023-01-29T16:31:59.05117+01:00"
    }
  ],
  "created_at": "2023-01-29T16:43:32.054728+01:00",
  "updated_at": "2023-01-29T16:43:32.054728+01:00"
}
```

Then, this role can be added to a principal using the `/v1/principals` endpoint.

### Delete a role

```bash
$ curl -s -X DELETE \
-H "Authorization: Bearer $ACCESS_TOKEN" \
http://localhost:8080/v1/roles/post-manager | jq
{
  "success": true
}
```

## User

You can create a new Authz user using the following API.

Note that a user is different from a principal: a user is just a user that can access the Authz API or frontend (event if a principal is also created automatically for it).

### Declare a new user

```bash
$ curl -s -X POST \
  -H 'Content-type: application/json' \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  -d '{"username": "john-cofee"}' \
  http://localhost:8080/v1/users | jq

{
  "username": "john-cofee",
  "password": "fDjtOEt_Fb",
  "created_at": "2023-01-29T16:51:30.391168+01:00",
  "updated_at": "2023-01-29T16:51:30.391168+01:00"
}
```

### Delete a user

```bash
$ curl -s -X DELETE \
-H "Authorization: Bearer $ACCESS_TOKEN" \
http://localhost:8080/v1/users/john-cofee | jq
{
  "success": true
}
```