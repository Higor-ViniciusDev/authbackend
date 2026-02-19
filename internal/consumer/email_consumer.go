package consumer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Higor-ViniciusDev/auth/configuration/logger"
	"github.com/Higor-ViniciusDev/auth/internal/entity"
	"github.com/Higor-ViniciusDev/auth/internal/internal_error"
	"github.com/Higor-ViniciusDev/auth/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

type EmailPendingConsumer struct {
	channel *amqp.Channel
}
type emailPendingMessage struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

func NewEmailPendingConsumer(channel *amqp.Channel) *EmailPendingConsumer {
	return &EmailPendingConsumer{
		channel: channel,
	}
}

// Start initializes the necessary infrastructure and begins consuming.
// Should be called in a separate goroutine in main.
func (c *EmailPendingConsumer) Start() {
	// Ensures the queue exists before attempting to consume
	if err := rabbitmq.SetupEmailPendingQueue(c.channel); err != nil {
		logger.Error("EmailPendingConsumer: error in declared email.pending", err)
		return
	}

	logger.Info("EmailPendingConsumer: waiting messages in email.pending")

	out := make(chan amqp.Delivery)

	// Starts listening in a separate goroutine to avoid blocking Start()
	go func() {
		if err := rabbitmq.Consumer(c.channel, "email.pending", out); err != nil {
			logger.Error("EmailPendingConsumer: error  consumer", err)
			close(out)
		}
	}()

	// Processing loop — blocks until the channel is closed
	for msg := range out {
		c.handle(msg)
	}

	logger.Info("EmailPendingConsumer: channel closed, consumer stopped")
}

// handle process the message received.
func (c *EmailPendingConsumer) handle(msg amqp.Delivery) {
	var payload emailPendingMessage

	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		logger.Error("EmailPendingConsumer: invalid payload", err)
		msg.Nack(false, false) // discard corrupted message
		return
	}

	if err := c.sendVerificationEmail(payload.Email, payload.Token); err != nil {
		logger.Error("EmailPendingConsumer: error send email", err)
		msg.Nack(false, true) // relocad in filla for try again
		return
	}

	msg.Ack(false) // confirms it was processed successfully
	logger.Info(fmt.Sprintf("EmailPendingConsumer: verification email sent to %s", payload.Email))
}

// sendVerificationEmail creates the link with the JWT and sends it via SMTP.
func (c *EmailPendingConsumer) sendVerificationEmail(email, token string) *internal_error.InternalError {
	url := os.Getenv("APP_URL")
	link := fmt.Sprintf("%s/verify-callback?token=%s", url, token)

	emailSenderMock := &entity.Email{
		To:      []string{email},
		Subject: "Confirm your registration",
		Body:    link,
	}

	if err := emailSenderMock.SendEmail(); err != nil {
		return internal_error.NewInternalServerError(err.Message)
	}

	return nil
}
