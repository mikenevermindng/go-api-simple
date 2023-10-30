package db

import (
	"context"
	_ "github.com/lib/pq"
	"goApiByGin/ent"
	"log"
)

var dbClient *ent.Client
var err error

func InitialDatabase(connectionString string) (*ent.Client, error) {
	dbClient, err = ent.Open(
		"postgres",
		connectionString,
	)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	} else {
		log.Println("connected to database")
	}

	if err := dbClient.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}
	return dbClient, err
}

func Client() *ent.Client {
	return dbClient
}
