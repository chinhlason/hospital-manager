package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

const (
	AcksRequireNone = kafka.RequireNone
	AcksRequireOne  = kafka.RequireOne
	AcksRequireAll  = kafka.RequireAll
)

const (
	_defaultKafkaSleep = 250 * time.Microsecond
	_defaultRetries    = 5
)

type ProducercConfig struct {
	Brokers      []string
	Topic        string
	BatchSize    int
	BatchTimeout time.Duration
	RequiredAcks kafka.RequiredAcks
}

type Producer[T Payload] struct {
	writer *kafka.Writer
}

func CreateProducer(cfg ProducercConfig) *Producer[Payload] {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(cfg.Brokers...),
		Topic:                  cfg.Topic,
		BatchSize:              cfg.BatchSize,
		BatchTimeout:           cfg.BatchTimeout,
		RequiredAcks:           cfg.RequiredAcks,
		AllowAutoTopicCreation: true,
	}
	return &Producer[Payload]{
		writer: writer,
	}
}

func (producer *Producer[T]) Produce(schema *Schema[T]) error {
	jsonData, err := json.Marshal(schema)
	if err != nil {
		return err
	}

	key, err := json.Marshal(schema.Key)
	if err != nil {
		return err
	}

	for i := 0; i < _defaultRetries; i++ {
		err = producer.writer.WriteMessages(context.Background(), kafka.Message{
			Value: jsonData,
			Key:   key,
		})
		if err != nil {
			log.Println("kafka-write-error", err)
			time.Sleep(_defaultKafkaSleep)
			continue
		} else {
			break
		}
	}
	return nil
}
