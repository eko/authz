Authz - Frontend
================

This is the frontend web UI of Authz.

Written in React, built with [Vite](https://vite.dev) and tested with [Vitest](https://vitest.dev).

## Pre-requisites

You will need the Authz backend server to run first. Please refer to backend [`README.md`](https://github.com/eko/authz/tree/master/backend) file.

## How to run

Install dependencies (Node.js 24+):

```bash
$ npm ci
```

Then run the development server (on http://localhost:3000) with:

```bash
$ npm run start
```

Or you can build static files (into `build/`) with:

```bash
$ npm run build
```

You can type-check the project with:

```bash
$ npm run typecheck
```

## Configuration

Here are the available configuration options available as environment variable **at build time only** (they are inlined by Vite, variables prefixed with `VITE_` are also exposed):

| Property | Default value | Description |
| -------- | ------------- | ----------- |
| REACT_APP_API_BASE_URI | `http://localhost:8080/v1` | Authz HTTP API backend url |
| REACT_APP_DEMO_ENABLED | N/A | Is demo mode enabled? (any non-empty value enables it) |
| REACT_APP_OAUTH_BUTTON_LABEL | `Sign-in wigh Single Sign-On (SSO)` | Sign in button label (will make it appear on front) |
| REACT_APP_OAUTH_ENABLED | `false` | Is OAuth sign-in enabled? |
| REACT_APP_OAUTH_LOGO_URL | N/A | Sign in button logo URL that will appear on left |

## Tests

### Unit tests

You can run unit tests (Vitest) with:

```bash
$ npm run test
```

Or in watch mode with:

```bash
$ npm run test:watch
```
