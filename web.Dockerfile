FROM golang:1.27.4 AS build
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /go/bin/kowalski ./cmd/web

FROM ghcr.io/greboid/dockerbase/nonroot:1.20261021.0
COPY --from=build /go/bin/kowalski /
CMD ["/kowalski"]
