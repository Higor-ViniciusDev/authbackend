package handlers

import (
	"encoding/json"
	"sync"

	"github.com/Higor-ViniciusDev/auth/configuration/logger"
	"github.com/Higor-ViniciusDev/auth/pkg/events"
	amqp "github.com/rabbitmq/amqp091-go"
)

// EmailVerifiedHandler publishes to the email.verified fanout exchange
type EmailVerifiedHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewEmailVerifiedHandler(rabbitMQChannel *amqp.Channel) *EmailVerifiedHandler {
	return &EmailVerifiedHandler{RabbitMQChannel: rabbitMQChannel}
}

func (h *EmailVerifiedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()

	jsonOutput, _ := json.Marshal(event.GetPayload())

	err := h.RabbitMQChannel.Publish(
		"email.verified",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonOutput,
		},
	)

	if err != nil {
		logger.Error("EmailVerifiedHandler: error publishing to exchange email.verified", err)
	}
}
