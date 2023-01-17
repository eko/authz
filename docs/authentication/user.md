# Authentication - Username / Password

Authz setups a default admin user with these credentials:

* Username: `admin`
* Default password: `changeme`

You can change this passwordb y using the `USER_ADMIN_DEFAULT_PASSWORD` configuration environment variable when launching the backend.

## Create a new user

When you will create a new user, a default password will be rendered, both when using the frontend or the HTTP API directly.

Creating a new user will also lead to the creation of a new `principal` with identifier prefixed by `authz-user-`.