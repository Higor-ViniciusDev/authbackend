package handlers

import (
	"encoding/json"
	"fmt"
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

	logger.Info(fmt.Sprintf("User criada: %v", event.GetPayload()))
	jsonOutput, _ := json.Marshal(event.GetPayload())

	msgRabbitmq := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	h.RabbitMQChannel.Publish(
		"amq.direct", // exchange
		"",           // key name
		false,        // mandatory
		false,        // immediate
		msgRabbitmq,  // message to publish
	)
}
