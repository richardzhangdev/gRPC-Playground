package producer

import (
	"context"
	"encoding/json"
	"gRPC-Playground/types"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	w := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{w}
}

func (p *Producer) PublishUsage(event types.UsageEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Value: data,
	}

	return p.writer.WriteMessages(context.Background(), msg)
}
