package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func connectRabbitmq(rabbitmqURL string) (*amqp.Channel, *amqp.Queue, func()) {
	conn, err := amqp.Dial(rabbitmqURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err := conn.Channel()
	failOnError(err, "Failedd to open a channel")

	q, err := ch.QueueDeclare(
		"hello",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	return ch, &q, func() {
		defer conn.Close()
		defer ch.Close()
	}
}

func main() {
	rabbitmqURL := "amqp://guest:guest@172.30.52.239:5672/"
	ch, q, closeConn := connectRabbitmq(rabbitmqURL)
	defer closeConn()

	body := "hello World!"
	err := ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		failOnError(err, "Failed to publish a message")
	} else {
		fmt.Println("Message published: %s", body)
	}

}
