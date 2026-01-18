package mqtt

import (
	"context"
	"fmt"
	"net/url"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
)

// Client wraps the autopaho connection manager
type Client struct {
	conn *autopaho.ConnectionManager
}

// NewClient creates a new MQTT client with auto-reconnect
func NewClient(ctx context.Context, brokerURL string, clientID string, onMessage func(topic string, payload []byte)) (*Client, error) {
	u, err := url.Parse(brokerURL)
	if err != nil {
		return nil, fmt.Errorf("invalid broker URL: %w", err)
	}

	cfg := autopaho.ClientConfig{
		ServerUrls:                    []*url.URL{u},
		KeepAlive:                     30,
		CleanStartOnInitialConnection: false,
		SessionExpiryInterval:         60,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			fmt.Println("[MQTT] Connection established")
			// Subscribe on connect (ensures subscription is reestablished on reconnect)
			if _, err := cm.Subscribe(ctx, &paho.Subscribe{
				Subscriptions: []paho.SubscribeOptions{
					{Topic: "device/1/command", QoS: 1},
				},
			}); err != nil {
				fmt.Printf("[MQTT] Failed to subscribe: %s\n", err)
			} else {
				fmt.Println("[MQTT] Subscribed to device/1/command")
			}
		},
		OnConnectError: func(err error) {
			fmt.Printf("[MQTT] Connection error: %s\n", err)
		},
		ClientConfig: paho.ClientConfig{
			ClientID: clientID,
			OnPublishReceived: []func(paho.PublishReceived) (bool, error){
				func(pr paho.PublishReceived) (bool, error) {
					onMessage(pr.Packet.Topic, pr.Packet.Payload)
					return true, nil
				},
			},
			OnClientError: func(err error) {
				fmt.Printf("[MQTT] Client error: %s\n", err)
			},
		},
	}

	conn, err := autopaho.NewConnection(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}

	// Wait for connection
	if err := conn.AwaitConnection(ctx); err != nil {
		return nil, fmt.Errorf("failed to await connection: %w", err)
	}

	return &Client{conn: conn}, nil
}

// Publish sends a message to the specified topic
func (c *Client) Publish(ctx context.Context, topic string, payload []byte, qos byte, retain bool) error {
	_, err := c.conn.Publish(ctx, &paho.Publish{
		Topic:   topic,
		QoS:     qos,
		Retain:  retain,
		Payload: payload,
	})
	return err
}

// Disconnect gracefully disconnects from the broker
func (c *Client) Disconnect(ctx context.Context) error {
	return c.conn.Disconnect(ctx)
}
