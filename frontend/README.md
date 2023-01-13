Authz - Frontend
================

This is the frontend web UI of Authz.

Written in React.

## Pre-requisites

You will need the Authz backend server to run first. Please refer to backend [`README.md`](https://github.com/eko/authz/tree/master/backend) file.

## How to run

You can simply run it with:

```bash
$ npm run start
```

Or you can build static files with:

```bash
$ npm run build
```

## Configuration

Here are the available configuration options available as environment variable **at build time only**:

| Property | Default value | Description |
| -------- | ------------- | ----------- |
| REACT_APP_API_BASE_URI | `http://localhost:8080/v1` | Authz HTTP API backend url |

## Tests

### Unit tests

You can run unit tests with:

```bash
$ npm run test
```
