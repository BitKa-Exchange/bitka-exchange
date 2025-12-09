package kafka

import (
	"bitka/services/account/internal/domain"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"time"
	"log"
	"strings"
)

type Consumer struct {
	uc domain.AccountUsecase
}

func NewConsumer(uc domain.AccountUsecase) *Consumer {
	return &Consumer{uc: uc}
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
		go c.handleMessages(partition)
		go c.handleErrors(partition)
		return
	}
}


func (c *Consumer) handleMessages(partition sarama.PartitionConsumer) {
	for msg := range partition.Messages() {
		if msg == nil {
			continue
		}

		var evt UserRegisteredEvent
		if err := json.Unmarshal(msg.Value, &evt); err != nil {
			log.Println("Failed to unmarshal Kafka message:", err)
			continue
		}

		if err := c.uc.CreateUserProfile(evt.UserID, evt.Email, evt.Username); err != nil {
			if isDuplicateError(err) {
				log.Println("Duplicate user profile, skipping:", evt.UserID)
			} else {
				log.Println("Failed to create user profile:", err)
			}
		} else {
			log.Println("User profile created:", evt.UserID)
		}
	}
}

func (c *Consumer) handleErrors(partition sarama.PartitionConsumer) {
	for err := range partition.Errors() {
		log.Println("Kafka partition error:", err)
	}
}

func isDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "duplicate key")
}

type UserRegisteredEvent struct {
	UserID   uuid.UUID `json:"user_id"`
	Email    string    `json:"email"`
	Username string    `json:"username"`
}
