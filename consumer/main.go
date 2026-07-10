package consumer

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"gRPC-Playground/metering"
	"gRPC-Playground/types"
)

func StartConsumer(brokers []string, topic string, groupID string, meterClient *metering.OpenMeterClient) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	defer r.Close()

	log.Println("Kafka consumer started")
	ctx := context.Background()

	for {
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			log.Printf("error reading message %v", err)
			continue
		}

		var usage types.UsageEvent
		err = json.Unmarshal(msg.Value, &usage)
		if err != nil {
			log.Printf("error unmarshalling message %v", err)
			continue
		}

		if err := meterClient.IngestEvent(ctx, &usage); err != nil {
			log.Printf("openmeter ingest failed %v", err)
		}
		//RecordUsage(usage)
	}
}

func RecordUsage(usage types.UsageEvent) {
	log.Printf(
		"request_id=%s model=%s prompt_tokens=%d completion_tokens=%d total_tokens=%d latency=%s",
		usage.RequestID,
		usage.Model,
		usage.PromptTokens,
		usage.CompletionTokens,
		usage.TotalTokens,
		usage.Latency,
	)
}
