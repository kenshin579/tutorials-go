FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY main.go .
RUN go build -o main-app main.go

FROM alpine:3.19
COPY --from=builder /app/main-app /usr/local/bin/
ENTRYPOINT ["main-app"]
