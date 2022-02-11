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
	RabbitMQ         struct {
		ConnectionString string    `json:"ConnectionString"`
		ChannelName      string    `json:"ChannelName"`
		Bindings         []Binding `json:"Bindings"`
	} `json:"RabbitMQ"`
}

func (s Settings) String() string {
	bytes, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		panic(fmt.Sprintln("Unable to json marshal Settings:", err))
	}

	return string(bytes)
}

type Binding struct {
	Exchange string `json:"Exchange"`
	Topic    string `json:"Topic"`
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
	s.ConnectionString = "user=postgres password=password dbname=rabbithole sslmode=verify-full"
	s.RabbitMQ.ConnectionString = "amqp://localhost/"
	s.RabbitMQ.ChannelName = "rabbithole"
	s.RabbitMQ.Bindings = make([]Binding, 0, 8)
}
