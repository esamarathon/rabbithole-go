package main

import (
	"encoding/json"
	"log"
	"os"
)

type FileRepository struct {
	file    *os.File
	encoder *json.Encoder
}

func CreateFileRepository(connectionString string) FileRepository {
	var file *os.File
	var err error
	switch connectionString {
	case "stdout":
		file = os.Stdout
	case "stderr":
		file = os.Stderr
	default:
		file, err = os.OpenFile(connectionString, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0733)
	}

	if err != nil {
		log.Fatalln("Unable to open file for writing: ", err)
	}

	return FileRepository{
		file:    file,
		encoder: json.NewEncoder(file),
	}
}

func (repo FileRepository) InsertEvent(e Event) error {
	return repo.encoder.Encode(e)
}
