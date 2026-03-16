package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	listenAddr := ":6379"
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis-service:6379"
	}

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", listenAddr, err)
	}
	defer listener.Close()
	log.Printf("redis-proxy listening on %s -> %s", listenAddr, redisAddr)

	for {
		client, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %v", err)
			continue
		}
		go handleConnection(client, redisAddr)
	}
}

func handleConnection(client net.Conn, redisAddr string) {
	defer client.Close()

	remote, err := net.Dial("tcp", redisAddr)
	if err != nil {
		log.Printf("failed to connect to redis at %s: %v", redisAddr, err)
		return
	}
	defer remote.Close()

	go io.Copy(remote, client)
	io.Copy(client, remote)
}
