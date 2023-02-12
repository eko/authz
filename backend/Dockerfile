FROM golang:1.20-alpine3.17 as builder

WORKDIR /project

COPY cmd/ ./cmd/
COPY configs/ ./configs/
COPY internal/ ./internal/
COPY pkg/authz/ ./pkg/authz/
COPY vendor/ ./vendor/
COPY go.mod .
COPY go.sum .

RUN apk add --no-cache gcc musl-dev

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=1 go build -mod vendor -ldflags "-s -w" -o authz ./cmd/

FROM alpine:3.17

RUN adduser -S app
USER app

COPY --from=builder /project/authz /authz

CMD [ "/authz" ]
