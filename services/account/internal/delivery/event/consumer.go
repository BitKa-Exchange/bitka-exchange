package event

import (
	"log"
	"time"

	"github.com/IBM/sarama"
)

type Consumer struct {
	handler *Handler
	brokers []string
	topic   string
}

func NewKafkaConsumer(handler *Handler) *Consumer {
	return &Consumer{
		handler: handler,
		brokers: []string{"kafka:9092"}, // default
		topic:   "user-registered",
	}
}

func (c *Consumer) Start() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	for {
		consumer, err := sarama.NewConsumer(c.brokers, config)
		if err != nil {
			log.Println("Kafka not ready, retrying in 2s:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		partition, err := consumer.ConsumePartition(c.topic, 0, sarama.OffsetNewest)
		if err != nil {
			log.Println("Kafka partition not ready, retrying in 2s:", err)
			time.Sleep(2 * time.Second)
			consumer.Close()
			continue
		}

		log.Println("Kafka consumer started ✓")
		go c.listenMessages(partition)
		go c.listenErrors(partition)

		return
	}
}

func (c *Consumer) listenMessages(partition sarama.PartitionConsumer) {
	for msg := range partition.Messages() {
		if msg == nil {
			continue
		}
		c.handler.HandleUserRegistered(msg.Value)
	}
}

func (c *Consumer) listenErrors(partition sarama.PartitionConsumer) {
	for err := range partition.Errors() {
		log.Println("Kafka error:", err)
	}
}
