package main

import (
	"log"

	"github.com/kenshin579/tutorials-go/go-grpc/chat"
	"golang.org/x/net/context"
)

type server struct {
	chat.UnimplementedChatServiceServer
}

func (s *server) SayHello(ctx context.Context, in *chat.Message) (*chat.Message, error) {
	log.Printf("Receive message body from client: %s", in.Body)
	return &chat.Message{Body: "Hello From the server!"}, nil
}
