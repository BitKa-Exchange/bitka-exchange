package event

import (
	"bitka/services/account/internal/domain"
)

type KafkaServer struct {
	consumer *Consumer
}

func NewKafkaServer(uc domain.AccountUsecase) *KafkaServer {
	server := StartKafkaServer(uc)
	go server.Start()
	return server
}

func StartKafkaServer(uc domain.AccountUsecase) *KafkaServer {
	handler := NewHandler(uc)
	consumer := NewKafkaConsumer(handler)
	return &KafkaServer{consumer: consumer}
}

func (s *KafkaServer) Start() {
	s.consumer.Start()
}
