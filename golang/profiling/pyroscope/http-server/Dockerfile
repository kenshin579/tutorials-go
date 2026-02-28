FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY http-server/ ./http-server/
RUN CGO_ENABLED=0 go build -o /app/http-server ./http-server/

FROM alpine:latest
COPY --from=builder /app/http-server /usr/local/bin/http-server
EXPOSE 8080
CMD ["http-server"]
