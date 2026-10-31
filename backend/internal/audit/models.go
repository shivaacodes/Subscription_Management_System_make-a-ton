package audit

import (
	"context"
	"encoding/json"
	"time"
)

type SubscriptionEvent struct {
	ID             string          `json:"id"`
	SubscriptionID string          `json:"subscriptionId"`
	EventType      string          `json:"eventType"`
	ActorID        *string         `json:"actorId,omitempty"`
	Payload        json.RawMessage `json:"payload"`
	CreatedAt      time.Time       `json:"createdAt"`
}

type Repository interface {
	RecordEvent(ctx context.Context, event *SubscriptionEvent) error
	ListEventsForSubscription(ctx context.Context, subscriptionID string) ([]SubscriptionEvent, error)
}
