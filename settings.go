package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Settings struct {
	Logging struct {
		Debug bool
	} `json:"Logging"`
	ConnectionString string
	RabbitMQ         RabbitSettings `json:"RabbitMQ"`
	Outputs []Output `json:"Outputs"`
}

func (s Settings) String() string {
	bytes, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		panic(fmt.Sprintln("Unable to json marshal Settings:", err))
	}

	return string(bytes)
}

func (settings *Settings) Migrate() {
	if settings.ConnectionString != "" {
		settings.Outputs = append(settings.Outputs, Output{
			Kind:             "sql",
			ConnectionString: settings.ConnectionString,
		})
		settings.ConnectionString = "" // Disable it now that we handled it.
	}
}

type RabbitSettings struct {
	ConnectionString string    `json:"ConnectionString"`
	ChannelName      string    `json:"ChannelName"`
	Bindings         []Binding `json:"Bindings"`
}

type Binding struct {
	Exchange string `json:"Exchange"`
	Topic    string `json:"Topic"`
}

type Output struct {
	Kind string
	ConnectionString string
}

func LoadSettings() (s Settings) {
	s.SetDefaults()
	file, err := os.Open("appsettings.json")
	if err != nil {
		log.Printf("Error reading config file: %s \n Progressing without it.", err)
		return
	}

	dec := json.NewDecoder(file)

	err = dec.Decode(&s)
	if err != nil {
		log.Fatalf("Error parsing config file: %s \n", err)
	}

	return
}

func (s *Settings) SetDefaults() {
	s.Logging.Debug = false
	s.ConnectionString = ""
	s.RabbitMQ.ConnectionString = "amqp://localhost/"
	s.RabbitMQ.ChannelName = "rabbithole"
	s.RabbitMQ.Bindings = []Binding{
		{
			Exchange: "demo",
			Topic:    "#",
		},
	}
	s.Outputs = []Output{
		{
			Kind: "sql",
			ConnectionString: "user=postgres password=password dbname=rabbithole sslmode=verify-full",
		},
	}
}
