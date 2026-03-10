package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
)

// InboundMessageEvent mirrors what the Node.js gateway publishes
type InboundMessageEvent struct {
	MessageID     string          `json:"message_id"`
	WabaID        string          `json:"waba_id"`
	PhoneNumberID string          `json:"phone_number_id"`
	From          string          `json:"from"`
	To            string          `json:"to"`
	Type          string          `json:"type"`
	Content       json.RawMessage `json:"content"`
	Timestamp     string          `json:"timestamp"`
	Contact       struct {
		WaID string `json:"wa_id"`
		Name string `json:"name"`
	} `json:"contact"`
}

type StatusUpdateEvent struct {
	MessageID   string `json:"message_id"`
	WabaID      string `json:"waba_id"`
	Status      string `json:"status"` // sent, delivered, read, failed
	Timestamp   string `json:"timestamp"`
	RecipientID string `json:"recipient_id"`
}

type Subscriber struct {
	redis    *redis.Client
	handler  *InboundHandler
	wabaIDs  []string // tenant WABA IDs to subscribe to
}

func NewSubscriber(rdb *redis.Client, handler *InboundHandler) *Subscriber {
	return &Subscriber{
		redis:   rdb,
		handler: handler,
	}
}

// Start subscribes to all tenant channels.
// In production, this should dynamically subscribe as new tenants are added.
func (s *Subscriber) Start(ctx context.Context) {
	// Subscribe to wildcard pattern for all WABAs
	// Pattern: usegro:inbound:message:*
	pubsub := s.redis.PSubscribe(ctx,
		"usegro:inbound:message:*",
		"usegro:inbound:status:*",
	)

	log.Println("Messaging subscriber listening for WhatsApp events...")

	go func() {
		ch := pubsub.Channel()
		for msg := range ch {
			switch {
			case matchPattern("usegro:inbound:message:*", msg.Channel):
				var event InboundMessageEvent
				if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
					log.Printf("failed to parse inbound message: %v", err)
					continue
				}
				if err := s.handler.HandleInbound(ctx, &event); err != nil {
					log.Printf("failed to handle inbound message %s: %v", event.MessageID, err)
				}

			case matchPattern("usegro:inbound:status:*", msg.Channel):
				var event StatusUpdateEvent
				if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
					log.Printf("failed to parse status update: %v", err)
					continue
				}
				if err := s.handler.HandleStatusUpdate(ctx, &event); err != nil {
					log.Printf("failed to handle status update %s: %v", event.MessageID, err)
				}
			}
		}
	}()

	// Also subscribe to outbound send requests from Go → Node.js
	// (other Go services publish here, this subscriber relays to the whatsapp-gateway via Redis)
}

// PublishOutbound publishes a send-message request for the Node.js gateway to pick up
func PublishOutbound(ctx context.Context, rdb *redis.Client, wabaID string, payload interface{}) error {
	channel := fmt.Sprintf("usegro:outbound:send:%s", wabaID)
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return rdb.Publish(ctx, channel, data).Err()
}

func matchPattern(pattern, channel string) bool {
	// Simple glob match for our fixed patterns
	if len(pattern) < 2 {
		return pattern == channel
	}
	prefix := pattern[:len(pattern)-1] // strip trailing *
	return len(channel) >= len(prefix) && channel[:len(prefix)] == prefix
}

func waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down messaging service...")
}
