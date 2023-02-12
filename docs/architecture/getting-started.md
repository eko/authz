# Getting started

We bring some tools to you so you can get started very easily.

## Docker

### Standalone image

The simple way to get started is to use our [`standalone`](https://github.com/eko/authz/blob/master/Dockerfile.standalone) Docker image which contains both backend and frontend in a single Docker image.

You can also use `SQLite` as database so you don't have to install any more dependency.

Just type:

```bash
$ docker run --rm \
    -e database_driver=sqlite \
    -e database_name=:memory: \
    -p 8080:8080 \
    -p 8081:8081 \
    -p 3000:80 \
    ekofr/authz:v0.8.0-standalone
```

Here, we are forwarding the 3 following ports:

* `8080`: the HTTP API,
* `8081`: the gRPC API,
* `3000`: the frontend UI.

Head to [`http://localhost:3000`](http://localhost:3000) on your machine and log in with following credentials:
* Username: `admin`
* Password: `changeme`

Now, you can play!

### Backend and Frontend images

If you prefer, you can use the separated backend and frontend Docker images:

```bash
$ docker run --rm \
    -e database_driver=sqlite \
    -e database_name=:memory: \
    -p 8080:8080 \
    -p 8081:8081 \
    ekofr/authz:v0.8.0-backend
```

and the frontend:

```bash
$ docker run --rm \
    -p 3000:80 \
    ekofr/authz:v0.8.0-frontend
```

## Build from sources

You can also clone the GitHub repository and build binaries and frontend static files from sources.

### Backend

```
$ git clone git@github.com:eko/authz.git
$ cd backend
$ go build -o backend ./cmd/main.go
```

Then, you can execute `./backend` to run the backend binary.

### Frontend

Static files can be generated with:

```bash
$ cd frontend
$ npm install
$ npm run build
```

Sources will then be available under `build/` directory.