package event

import (
	"github.com/IBM/sarama"
	"log"
	"time"
)

type Consumer struct {
	handler *Handler
}

func NewConsumer(handler *Handler) *Consumer {
	return &Consumer{handler: handler}
}

func (c *Consumer) Start() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	for {
		consumer, err := sarama.NewConsumer([]string{"kafka:9092"}, config)
		if err != nil {
			log.Println("Kafka not ready, retrying in 2s:", err)
			time.Sleep(2 * time.Second)
			continue
		}

		partition, err := consumer.ConsumePartition("user-registered", 0, sarama.OffsetNewest)
		if err != nil {
			log.Println("Partition consumer not ready, retrying in 2s:", err)
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
		log.Println("Kafka partition error:", err)
	}
}
