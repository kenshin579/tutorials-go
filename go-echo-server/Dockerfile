FROM golang:alpine AS build
WORKDIR /app
ADD . /app

RUN go build -o bin/echo-server cmd/echo-server/main.go

FROM alpine

WORKDIR /app
COPY --from=build /app/bin/echo-server /app

CMD ["/app/echo-server"]
