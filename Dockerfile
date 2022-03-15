FROM reg.c5h.io/golang AS build
WORKDIR /app
COPY . /app
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /go/bin/kowalski ./cmd/kowalski

FROM reg.c5h.io/base
COPY --from=build /go/bin/kowalski /
CMD ["/kowalski"]
