package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"time"
)

const (
	ConsumeAllData     = kafka.FirstOffset
	ConsumeNewDataOnly = kafka.LastOffset
)

type ConsumerConfig struct {
	Brokers     []string
	Topic       string
	GroupId     string
	MinBytes    int
	MaxBytes    int
	MaxWait     time.Duration
	StartOffSet int64
}

type Consume[T Payload] struct {
	context context.Context
	reader  *kafka.Reader
}

func CreateConsumer(cfg ConsumerConfig) *Consume[Payload] {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     cfg.Brokers,
		Topic:       cfg.Topic,
		GroupID:     cfg.GroupId,
		MinBytes:    cfg.MinBytes,
		MaxBytes:    cfg.MaxBytes,
		MaxWait:     cfg.MaxWait,
		StartOffset: cfg.StartOffSet,
	})
	return &Consume[Payload]{
		context: context.Background(),
		reader:  reader,
	}
}

func (consumer *Consume[T]) Consume() (*Schema[T], error) {
	data, err := consumer.reader.ReadMessage(consumer.context)
	if err != nil {
		return nil, err
	}
	var schema = &Schema[T]{}
	err = json.Unmarshal(data.Value, &schema)
	if err != nil {
		return nil, err
	}
	return schema, err
}
