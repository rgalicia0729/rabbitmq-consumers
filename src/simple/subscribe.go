package simple

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"time"
)

func Subscribe() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Panicf("%s: %s", err, "Failed to connect to RabbitMQ")
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("%s: %s", err, "Failed to open a channel")
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"q.simple.event",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Panicf("%s: %s", err, "Failed to declare a queue")
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Panicf("%s: %s", err, "Failed to set QoS")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Panicf("%s: %s", err, "Failed to register a consumer")
	}

	var forever chan struct{}
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			time.Sleep(time.Duration(5 * time.Second))
			log.Println("Message processed!!")
			d.Ack(false)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
