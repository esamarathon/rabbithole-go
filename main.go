package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

func main() {
	log.Println("Starting Rabbithole")
	defer func() {
		log.Println("Shutting down.")
	}()

	settings := LoadSettings()
	settings.Migrate()


	log.Println("Settings: ")
	log.Println(settings)

	amqpConn := Connect(settings)

	wg := new(sync.WaitGroup)

	for _, output := range settings.Outputs {
		wg.Add(1)
		go HandleOutput(output, amqpConn, settings.RabbitMQ, wg)
	}

	wg.Wait()
	log.Println("End of message queue.")
}



func HandleOutput(output Output, amqpConn *amqp.Connection, settings RabbitSettings, wg *sync.WaitGroup) {
	repo := MakeRepository(output)
	log.Printf("Setup complete. Starting message handling loop for %s, \n", output.Kind)
	stream := Listen(amqpConn, settings)
	ListenAndServe(stream, repo)
	wg.Done()
}


type Repository interface {
	InsertEvent(e Event) error
}

func MakeRepository(output Output) Repository {
	var repo Repository
	switch output.Kind {
	case "sql": repo = ConnectToDatabase(output.ConnectionString)
	case "file": repo = CreateFileRepository(output.ConnectionString)
	default: panic("Unknown output kind. Valid kinds: [sql,file]")
	}

	return repo;
}

func ListenAndServe(stream RabbitStream, repo Repository) {
	for msg := range stream.messages {
		log.Println("Incoming message.")
		content := make(map[string]interface{})
		err := json.Unmarshal(msg.Body, &content)
		if err != nil {
			//Failed to unserialize the event. So the message is busted and needs to be discarded.
			log.Println("Discarding invalid message from AMQP.")
			msg.Nack(false, false)
		}

		event := Event{
			Id:         uuid.New(),
			Recieved:   time.Now().UTC(),
			Exchange:   msg.Exchange,
			RoutingKey: msg.RoutingKey,
			Content:    content,
		}

		log.Printf("Storing message from %s (%s): %s", event.Exchange, event.RoutingKey, string(msg.Body))
		repo.InsertEvent(event)
		msg.Ack(false)
	}
}
