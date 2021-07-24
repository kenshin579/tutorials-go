FROM golang:alpine AS build
ADD . /buildspace
WORKDIR /buildspace
RUN go build -o helloworld ./...

FROM alpine
WORKDIR /app
COPY --from=build /buildspace/helloworld /app
CMD [ "/app/helloworld" ]
