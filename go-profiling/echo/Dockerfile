FROM golang:alpine AS build
WORKDIR /app
ADD . /app

RUN go build -o bin/profile-server main.go

FROM alpine

WORKDIR /app
COPY --from=build /app/bin/profile-server /app

CMD ["/app/profile-server"]
