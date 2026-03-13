FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /pull-server ./pull-server/

FROM alpine:latest
COPY --from=builder /pull-server /pull-server
EXPOSE 6060
CMD ["/pull-server"]
