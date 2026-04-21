package main

import (
	"log"

	"github.com/streadway/amqp"
)

func Listen(conn *amqp.Connection, settings RabbitSettings) RabbitStream {
	ch, q := PrepareListeningQueue(conn, settings)
	for _, b := range settings.Bindings {
		DeclareExchange(ch, b)
		ch.QueueBind(q.Name, b.Topic, b.Exchange, false, nil)
	}

	msgs, _ := ch.Consume(
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

func DeclareExchange(ch *amqp.Channel, b Binding) {
	// name
	// durable
	// delete when unused
	// exclusive
	// no-wait
	// arguments
	err := ch.ExchangeDeclare(b.Exchange, amqp.ExchangeTopic, true, true, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare exchange. ", err)
	}
}

func Connect(settings Settings) *amqp.Connection {
	conn, err := amqp.Dial(settings.RabbitMQ.ConnectionString)
	if err != nil {
		log.Fatal(err, "Failed to open connection to amqp server.")
	}
	return conn
}

func PrepareListeningQueue(conn *amqp.Connection, settings RabbitSettings) (*amqp.Channel, amqp.Queue) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err, "Failed to open channel.")
	}

	q, err := ch.QueueDeclare(
		settings.ChannelName,
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
