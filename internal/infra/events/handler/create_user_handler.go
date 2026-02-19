package handlers

import (
	"encoding/json"
	"sync"

	"github.com/Higor-ViniciusDev/auth/configuration/logger"
	"github.com/Higor-ViniciusDev/auth/pkg/events"
	amqp "github.com/rabbitmq/amqp091-go"
)

type UserCreateHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewCreateuserHandler(rabbitMQChannel *amqp.Channel) *UserCreateHandler {
	return &UserCreateHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *UserCreateHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()

	jsonOutput, _ := json.Marshal(event.GetPayload())

	msgRabbitmq := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	err := h.RabbitMQChannel.Publish(
		"",
		"email.pending",
		false,
		false,
		msgRabbitmq,
	)

	if err != nil {
		logger.Error("error publishing event to RabbitMQ", err)
	}
}
