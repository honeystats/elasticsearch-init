FROM golang:1.17-alpine as builder
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o /init

FROM alpine:latest
WORKDIR /
COPY --from=builder /init /init
COPY ./dashboards/*.ndjson /dashboards/
CMD ["/init"]
