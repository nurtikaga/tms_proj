package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/nurtikaga/tms_proj/internal/core/entity"
)

type eventMessage struct {
	ID         string    `json:"id"`
	ShipmentID string    `json:"shipment_id"`
	Status     string    `json:"status"`
	Note       string    `json:"note"`
	OccurredAt time.Time `json:"occurred_at"`
}

type Producer struct {
	writer *kafka.Writer
}

func New(brokerAddr, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokerAddr),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, e *entity.ShipmentEvent) error {
	msg := eventMessage{
		ID:         e.ID,
		ShipmentID: e.ShipmentID,
		Status:     string(e.Status),
		Note:       e.Note,
		OccurredAt: e.OccurredAt,
	}
	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("kafka: marshal event: %w", err)
	}
	if err = p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(e.ShipmentID),
		Value: payload,
	}); err != nil {
		return fmt.Errorf("kafka: write message: %w", err)
	}
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
