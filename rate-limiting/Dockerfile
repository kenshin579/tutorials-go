FROM golang:1.26-alpine AS builder

WORKDIR /app
COPY ../../go.mod ../../go.sum ../../
COPY . ./rate-limiting/

WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 go build -o /server ./rate-limiting/cmd/server

FROM alpine:3.19
COPY --from=builder /server /server
EXPOSE 8080
CMD ["/server"]
