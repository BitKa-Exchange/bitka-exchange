package event

import (
	"bitka/services/account/internal/delivery/event/dto"
	"bitka/services/account/internal/domain"
	"encoding/json"
	"log"
	"strings"
)

type Handler struct {
	uc domain.AccountUsecase
}

func NewHandler(uc domain.AccountUsecase) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) HandleUserRegistered(msg []byte) {
	var evt dto.UserRegisteredEvent
	if err := json.Unmarshal(msg, &evt); err != nil {
		log.Println("Failed to unmarshal Kafka message:", err)
		return
	}

	if err := h.uc.CreateUserProfile(evt.UserID, evt.Email, evt.Username); err != nil {
		if isDuplicateError(err) {
			log.Println("Duplicate user profile, skipping:", evt.UserID)
		} else {
			log.Println("Failed to create user profile:", err)
		}
	} else {
		log.Println("User profile created:", evt.UserID)
	}
}

func isDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "duplicate key")
}
