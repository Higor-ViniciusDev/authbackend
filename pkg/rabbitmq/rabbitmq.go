package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func OpenChannel() (*amqp.Channel, error) {
	amqpConn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		return nil, err
	}

	channel, err := amqpConn.Channel()
	if err != nil {
		return nil, err
	}

	return channel, nil
}

func SetupEmailPendingQueue(ch *amqp.Channel) error {
	_, err := ch.QueueDeclare(
		"email.pending",
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-message-ttl": int32(43200000), // TTL 12h em ms
		},
	)
	return err
}

// SetupEmailVerifiedExchange declares the email.verified exchange as fanout (idempotent).
func SetupEmailVerifiedExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"email.verified",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
}

// SetupPodQueue creates an exclusive queue for this pod and binds it to email.verified exchange.
func SetupPodQueue(ch *amqp.Channel) (string, error) {
	// empty name = RabbitMQ generates a unique name per pod
	q, err := ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return "", err
	}

	// binds the queue to the fanout exchange
	err = ch.QueueBind(
		q.Name,
		"",
		"email.verified",
		false,
		nil,
	)
	if err != nil {
		return "", err
	}

	return q.Name, nil
}

// Simples consumer
func Consumer(ch *amqp.Channel, queueName string, out chan amqp.Delivery) error {
	msgs, err := ch.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for msg := range msgs {
		out <- msg
	}

	return nil
}

func Publish(ch *amqp.Channel, body string, exName string) error {
	return ch.Publish(
		exName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
}
