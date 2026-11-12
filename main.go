package main

import (
	"log"
	"os"

	"github.com/aleksandradidenko/go_final_project/pkg/db"
	"github.com/aleksandradidenko/go_final_project/pkg/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	port := os.Getenv("TODO_PORT")
	dbFile := os.Getenv("TODO_DBFILE")
	password := os.Getenv("TODO_PASSWORD")

	if port == "" {
		port = "7540"
	}
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Ошибка закрытия базы данных: %v", err)
		}
	}()

	return server.Run(port, password)
}
