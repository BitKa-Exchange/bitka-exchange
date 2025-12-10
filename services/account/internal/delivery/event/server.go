package event

type Server struct {
	consumer *Consumer
}

func NewServer(handler *Handler) *Server {
	consumer := NewConsumer(handler)
	return &Server{consumer: consumer}
}

func (s *Server) Start() {
	s.consumer.Start()
}
