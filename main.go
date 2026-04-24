package main

import (
	"encoding/json"

	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

var (
	settings Settings
)

func main() {
	DebugLogln("Starting Rabbithole")
	defer func() {
		DebugLogln("Shutting down.")
	}()

	settings = LoadSettings()

	DebugLogln("Settings: ")
	DebugLogln(settings)

	events := make(chan Event)
	wg := new(sync.WaitGroup)

	outputChannels := make([]chan<- Event, 0, len(settings.Outputs))
	for _, output := range settings.Outputs {
		ch := make(chan Event)
		wg.Add(1)
		go HandleOutput(output, ch, wg)
		outputChannels = append(outputChannels, ch)
	}
	go FanOut(events, outputChannels...)

	amqpConn := Connect(settings)
	stream := Listen(amqpConn, settings.RabbitMQ)
	go Transform(stream.messages, events)
	
	wg.Wait()
	stream.Close()
	DebugLogln("End of message queue.")
}

func Transform(messages <-chan amqp.Delivery, events chan<- Event) {
	defer close(events)

	for msg := range messages {
		DebugLogln("Incoming message.")
		var content any
		content = make(map[string]interface{})
		err := json.Unmarshal(msg.Body, &content)
		if err != nil {
			content = string(msg.Body) //Not a JSON object body, so just dump value.
		}

		event := Event{
			Id:         uuid.New(),
			Received:   time.Now().UTC(),
			Exchange:   msg.Exchange,
			RoutingKey: msg.RoutingKey,
			Content:    content,
		}

		events <- event
		msg.Ack(false)
	}
}

func HandleOutput(output Output, events <-chan Event, wg *sync.WaitGroup) {
	defer wg.Done()

	repo := MakeRepository(output)
	DebugLogf("Setup complete. Starting message handling loop for %s, \n", output.Kind)

	for event := range events {
		DebugLogf("Storing message from %s (%s)", event.Exchange, event.RoutingKey)
		repo.InsertEvent(event)
	}
}

func FanOut[T interface{}](input <-chan T, outputs ...chan<- T) {
	defer func() {
		for _, ch := range outputs {
			close(ch)
		}
	}()

	for val := range input {
		for _, ch := range outputs {
			ch <- val
		}
	}
}
