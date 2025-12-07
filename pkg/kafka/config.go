package kafka

import (
	"app-ecommerce/config"
	"context"
	"log"
	"sync"

	"github.com/IBM/sarama"
)

type ProducerConfig struct {
	sarama.SyncProducer
}

var producerInsatnce *ProducerConfig
var producerOnce sync.Once

func InitProducer() {
	producerOnce.Do(func() {
		cfg := config.GetConfig()
		producer, err := sarama.NewSyncProducer(cfg.Kafka.Brokers, nil)
		if err != nil {
			log.Fatalf("Failed to connect to kafka producer: %v", err)
		}
		producerInsatnce = &ProducerConfig{producer}
	})
}

func CloseProducer() {
	producerInsatnce.Close()
}

type ConsumerConfig struct {
	sarama.ConsumerGroup
}

var consumerInsatnce *ConsumerConfig
var consumerOnce sync.Once

func InitConsumer() {
	consumerOnce.Do(func() {
		cfg := config.GetConfig()
		consumer, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, cfg.Kafka.GroupID, nil)
		if err != nil {
			log.Fatalf("Failed to connect to kafka consumer: %v", err)
		}
		consumerInsatnce = &ConsumerConfig{consumer}
	})
}

func GetConsumer() *ConsumerConfig {
	return consumerInsatnce
}

func (c *ConsumerConfig) RunConsumer(ctx context.Context, topics []string, consumer sarama.ConsumerGroupHandler) {
	c.Consume(ctx, topics, consumer)
}

func CloseConsumer() {
	consumerInsatnce.Close()
}
