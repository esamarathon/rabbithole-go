package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Repository struct {
	db              *sqlx.DB
	eventInsertStmt *sqlx.Stmt
}

var schema = `
CREATE TABLE IF NOT EXISTS public."Events" (
    id uuid NOT NULL PRIMARY KEY,
    "timestamp" timestamp without time zone NOT NULL,
    exchange text,
    routing_key text,
    content jsonb
);
`

type Event struct {
	Id         uuid.UUID              `db:"id"`
	Recieved   time.Time              `db:"timestamp"`
	Exchange   string                 `db:"exchange"`
	RoutingKey string                 `db:"routing_key"`
	Content    map[string]interface{} `db:"content"`
}

func ConnectToDatabase(settings Settings) (repo Repository) {
	var err error
	repo.db, err = sqlx.Connect("postgres", settings.ConnectionString)
	if err != nil {
		log.Fatalln("Unable to connect to database: ", err)
	}

	repo.db.Exec(schema) //Make sure DB is setup.

	repo.eventInsertStmt, err = repo.db.Preparex(`
	INSERT INTO 
		public."Events" (id, timestamp, exchange, routing_key, content) 
	VALUES 
		($1, $2, $3, $4, $5)`)
	if err != nil {
		log.Fatalln("Unable to prepare statement to insert Events: ", err)
	}

	return
}

func (repo Repository) InsertEvent(e Event) error {
	jsoncontent, err := json.Marshal(e.Content)
	if err != nil {
		return err
	}

	_, err = repo.eventInsertStmt.Exec(e.Id.String(), e.Recieved.Format(time.RFC3339), e.Exchange, e.RoutingKey, jsoncontent)
	return err
}
