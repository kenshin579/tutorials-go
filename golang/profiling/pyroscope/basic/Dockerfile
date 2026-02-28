FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY basic/ ./basic/
RUN CGO_ENABLED=0 go build -o /app/basic ./basic/

FROM alpine:latest
COPY --from=builder /app/basic /usr/local/bin/basic
CMD ["basic"]
