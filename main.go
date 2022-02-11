package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

func main() {
	log.Println("Starting CrowdControl Coin bridge.")
	defer func() {
		log.Println("Shutting down.")
	}()

	settings := LoadSettings()

	log.Println("Settings: ")
	log.Println(settings)

	repo := ConnectToDatabase(settings)
	stream := ConnectToRabbitMQ(settings)

	log.Println("Setup complete. Starting message handling loop.")
	ListenAndServe(stream, repo)
	log.Println("End of message queue.")
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
