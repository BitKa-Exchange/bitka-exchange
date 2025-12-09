package kafka

import (
	"encoding/json"
	"log"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type Producer struct {
	client sarama.SyncProducer
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{client: producer}, nil
}

type UserRegisterEvent struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
}

// PublishUserRegister sends the event to the Kafka topic
func (p *Producer) PublishUserRegister(event UserRegisterEvent) error {
	regis, err := json.Marshal(event)
	if err != nil {
		return err
	}
	message := &sarama.ProducerMessage{
		Topic: "user-registered",
		Value: sarama.ByteEncoder(regis),
	}

	partition, offset, err := p.client.SendMessage(message)
	if err != nil {
		return err
	}

	log.Printf("Message published to partition %d at offset %d\n", partition, offset)
	return nil
}
