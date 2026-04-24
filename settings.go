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
	RabbitMQ RabbitSettings `json:"RabbitMQ"`
	Outputs  []Output       `json:"Outputs"`
}

func (s Settings) String() string {
	bytes, err := json.MarshalIndent(s, "", "\t")
	if err != nil {
		panic(fmt.Sprintln("Unable to json marshal Settings:", err))
	}

	return string(bytes)
}

type RabbitSettings struct {
	ConnectionString string    `json:"ConnectionString"`
	ChannelName      string    `json:"ChannelName"`
	Bindings         []Binding `json:"Bindings"`
}

type Binding struct {
	Exchange Exchange `json:"Exchange"`
	Topic    string   `json:"Topic"`
}

type Exchange struct {
	Name       string                 `json:"Name"`
	Durable    bool                   `json:"Durable"`
	AutoDelete bool                   `json:"AutoDelete"`
	Arguments  map[string]interface{} `json:"Arguments"`
}

func (e *Exchange) UnmarshalJSON(data []byte) error {
	type ExchangeAlias Exchange
	test := &ExchangeAlias{
		Durable:    true,
		AutoDelete: true,
	}

	err := json.Unmarshal(data, test)
	if err != nil {
		err2 := json.Unmarshal(data, &test.Name)
		if err2 != nil {
			return err //Yes, return err. err2 is only to see if this fallback method worked. If that also failed, we want original error.
		}
	}

	*e = Exchange(*test)
	return nil
}

type Output struct {
	Kind             string
	ConnectionString string
}

func LoadSettings() (s Settings) {
	s.SetDefaults()
	file, err := os.Open("appsettings.json")
	if err != nil {
		Logf("Error reading config file: %s \n Progressing without it.", err)
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
	s.RabbitMQ.ConnectionString = "amqp://localhost/"
	s.RabbitMQ.ChannelName = "rabbithole"
	s.RabbitMQ.Bindings = []Binding{
		{
			Exchange: Exchange{
				Name: "demo",
			},
			Topic: "#",
		},
	}
	s.Outputs = []Output{
		{
			Kind:             "sql",
			ConnectionString: "user=postgres password=password dbname=rabbithole sslmode=verify-full",
		},
	}
}
