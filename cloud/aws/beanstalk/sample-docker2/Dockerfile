FROM golang:1.19.5-alpine3.17 AS build-env
WORKDIR /app

COPY . .

RUN go build -o application .

FROM alpine

EXPOSE 8000

# Run it
ENTRYPOINT ["./application"]
