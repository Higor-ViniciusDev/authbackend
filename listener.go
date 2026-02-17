package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	// Conecta ao RabbitMQ (mesma URL usada no pkg/rabbitmq/rabbitmq.go)
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Declara uma exchange (opcional, se já existir, mas garante que amq.direct está lá)
	err = ch.ExchangeDeclare(
		"amq.direct", // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	// Declara uma fila temporária exclusiva
	q, err := ch.QueueDeclare(
		"",    // name (vazio para gerar um nome aleatório)
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Faz o bind da fila com a exchange
	// Como o publisher usa routing key vazia no create_user_handler.go line 36: "",
	// nós fazemos o bind com routing key vazia também.
	log.Printf("Binding queue %s to exchange amq.direct with routing key ''", q.Name)
	err = ch.QueueBind(
		q.Name,       // queue name
		"",           // routing key
		"amq.direct", // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [x] Mensagem Recebida: %s", d.Body)
		}
	}()

	log.Printf(" [*] Aguardando mensagens em amq.direct. Para sair pressione CTRL+C")
	<-forever
}
