package main

import (
	"log"

	"github.com/streadway/amqp"
)

func ConnectToRabbitMQ(settings Settings) RabbitStream {
	conn, err := amqp.Dial(settings.RabbitMQ.ConnectionString)
	if err != nil {
		log.Fatal(err, "Failed to open connection to amqp server.")
	}

	// name
	// durable
	// delete when unused
	// exclusive
	// no-wait
	// arguments
	ch, q := PrepareListeningQueue(conn, settings)

	for _, b := range settings.RabbitMQ.Bindings {
		err = ch.ExchangeDeclare(b.Exchange, amqp.ExchangeTopic, true, true, false, false, nil)
		if err != nil {
			log.Fatal("Failed to declare exchange. ", err)
		}

		if err != nil {
			log.Fatal("Unable to declare exchange")
		}

		ch.QueueBind(q.Name, b.Topic, b.Exchange, false, nil)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	return RabbitStream{
		conn:     conn,
		channel:  ch,
		messages: msgs,
	}
}

func PrepareListeningQueue(conn *amqp.Connection, settings Settings) (*amqp.Channel, amqp.Queue) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err, "Failed to open channel.")
	}

	q, err := ch.QueueDeclare(
		settings.RabbitMQ.ChannelName,
		true,
		false,
		false,
		false,
		nil,
	)
	return ch, q
}

type RabbitStream struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	messages <-chan amqp.Delivery
}

func (s RabbitStream) Close() {
	s.channel.Close()
	s.conn.Close()
}
