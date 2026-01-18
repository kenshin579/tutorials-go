package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-mqtt-dashboard/backend/internal/device"
	"go-mqtt-dashboard/backend/internal/mqtt"
)

const (
	brokerURL     = "mqtt://localhost:1883"
	clientID      = "go-backend-device-1"
	stateTopic    = "device/1/state"
	publishPeriod = 2 * time.Second
)

func main() {
	// Setup context with signal handling
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create device simulator
	sim := device.NewSimulator()

	// Message handler for commands
	onMessage := func(topic string, payload []byte) {
		var cmd struct {
			Action string `json:"action"`
		}
		if err := json.Unmarshal(payload, &cmd); err != nil {
			log.Printf("[CMD] Failed to parse command: %s", err)
			return
		}
		sim.HandleCommand(cmd.Action)
		log.Printf("[CMD] Command received: %s", cmd.Action)
	}

	// Connect to MQTT broker
	log.Println("[MAIN] Connecting to MQTT broker...")
	client, err := mqtt.NewClient(ctx, brokerURL, clientID, onMessage)
	if err != nil {
		log.Fatalf("[MAIN] Failed to connect: %s", err)
	}
	log.Println("[MAIN] Connected to MQTT broker")

	// Publish state every 2 seconds
	ticker := time.NewTicker(publishPeriod)
	defer ticker.Stop()

	log.Println("[MAIN] Starting state publish loop...")
	for {
		select {
		case <-ctx.Done():
			log.Println("[MAIN] Shutting down...")
			client.Disconnect(context.Background())
			return
		case <-ticker.C:
			state := sim.GetState()
			payload, err := json.Marshal(state)
			if err != nil {
				log.Printf("[STATE] Failed to marshal state: %s", err)
				continue
			}

			if err := client.Publish(ctx, stateTopic, payload, 0, true); err != nil {
				log.Printf("[STATE] Failed to publish: %s", err)
				continue
			}

			fmt.Printf("[STATE] Published: status=%s, temperature=%.1f°C\n", state.Status, state.Temperature)
		}
	}
}
