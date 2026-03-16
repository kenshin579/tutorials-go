FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY main.go .
RUN go build -o metrics-adapter main.go

FROM alpine:3.19
COPY --from=builder /app/metrics-adapter /usr/local/bin/
ENTRYPOINT ["metrics-adapter"]
