package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Message is the type of messages from the client
type Message struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}

func (m Message) String() string {
	return fmt.Sprintf("%s: %s", m.Sender, m.Text)
}

// Response is the type of response to the client
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (r Response) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		log.Println("could not marshal response:", err)
		return ""
	}
	return string(b)
}

// Broker accepts subscriptions from clients and publishes messages to them all
type Broker struct {
	subscribers map[chan Message]bool
	sync.Mutex
}

// Subscribe adds a client to the broker
func (b *Broker) Subscribe() chan Message {
	b.Lock()
	defer b.Unlock()
	log.Println("subscribing to broker")
	ch := make(chan Message)
	b.subscribers[ch] = true
	return ch
}

// Unsubscribe removes a client from the broker
func (b *Broker) Unsubscribe(ch chan Message) {
	b.Lock()
	defer b.Unlock()
	log.Println("unsubscribing from broker")
	close(ch)
	delete(b.subscribers, ch)
}

// Publish sends a slice of bytes to all subscribed clients
func (b *Broker) Publish(msg Message) {
	b.Lock()
	defer b.Unlock()
	log.Printf("Publishing to %d subscribers\n", len(b.subscribers))
	for ch := range b.subscribers {
		ch <- msg
	}
}

// NewBroker creates a new broker
func NewBroker() *Broker {
	return &Broker{subscribers: make(map[chan Message]bool)}
}

func messageHandler(b *Broker) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if r.Header.Get("Content-Type")[:16] != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusNotAcceptable)
			return
		}

		var m Message
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&m); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b.Publish(m)

		if r.Header.Get("Accept") == "application/json" {
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, Response{Success: true, Message: "OK"})
			return
		}

		fmt.Fprintln(w, "OK")
	}
}

func (b *Broker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	f, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := b.Subscribe()
	defer b.Unsubscribe(ch)

	for {
		select {
		case msg := <-ch:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			f.Flush()
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	msgBroker := NewBroker()
	http.Handle("/events", msgBroker)
	http.HandleFunc("/update", messageHandler(msgBroker))
	http.Handle("/", http.FileServer(http.Dir(".")))
	err := http.ListenAndServe("localhost:8888", nil)
	if err != nil {
		panic(err)
	}
}
