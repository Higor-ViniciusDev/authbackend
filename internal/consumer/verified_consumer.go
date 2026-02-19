package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/Higor-ViniciusDev/auth/configuration/logger"
	"github.com/Higor-ViniciusDev/auth/internal/infra/ws"
	"github.com/Higor-ViniciusDev/auth/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

// EmailVerifiedConsumer consumes the email.verified fanout exchange.
type EmailVerifiedConsumer struct {
	channel   *amqp.Channel
	wsManager *ws.Manager
}

type emailVerifiedMessage struct {
	Email string `json:"email"`
}

func NewEmailVerifiedConsumer(channel *amqp.Channel, wsManager *ws.Manager) *EmailVerifiedConsumer {
	return &EmailVerifiedConsumer{
		channel:   channel,
		wsManager: wsManager,
	}
}

// Start sets up the exchange, creates this pod's exclusive queue and starts consuming.
func (c *EmailVerifiedConsumer) Start() {
	// 1. makes sure the fanout exchange exists
	if err := rabbitmq.SetupEmailVerifiedExchange(c.channel); err != nil {
		logger.Error("EmailVerifiedConsumer: error declaring exchange email.verified", err)
		return
	}

	// 2. creates exclusive queue for this pod and binds to exchange
	queueName, err := rabbitmq.SetupPodQueue(c.channel)
	if err != nil {
		logger.Error("EmailVerifiedConsumer: error creating pod queue", err)
		return
	}

	logger.Info(fmt.Sprintf("EmailVerifiedConsumer: pod queue created → %s", queueName))

	out := make(chan amqp.Delivery)

	go func() {
		if err := rabbitmq.Consumer(c.channel, queueName, out); err != nil {
			logger.Error("EmailVerifiedConsumer: error starting consumer", err)
			close(out)
		}
	}()

	for msg := range out {
		c.handle(msg)
	}

	logger.Info("EmailVerifiedConsumer: channel closed, consumer stopped")
}

func (c *EmailVerifiedConsumer) handle(msg amqp.Delivery) {
	var payload emailVerifiedMessage

	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		logger.Error("EmailVerifiedConsumer: invalid payload", err)
		msg.Nack(false, false)
		return
	}

	// notifies WebSocket on this pod for the received email.
	c.wsManager.Notify(payload.Email)

	msg.Ack(false)
	logger.Info(fmt.Sprintf("EmailVerifiedConsumer: notification processed for %s", payload.Email))
}
