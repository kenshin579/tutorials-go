FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY main.go .
RUN go build -o redis-proxy main.go

FROM alpine:3.19
COPY --from=builder /app/redis-proxy /usr/local/bin/
ENTRYPOINT ["redis-proxy"]
