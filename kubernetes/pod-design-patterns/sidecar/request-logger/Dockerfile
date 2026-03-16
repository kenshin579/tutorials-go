FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY main.go .
RUN go build -o request-logger main.go

FROM alpine:3.19
COPY --from=builder /app/request-logger /usr/local/bin/
ENTRYPOINT ["request-logger"]
