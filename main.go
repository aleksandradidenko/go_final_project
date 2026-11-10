package main

import (
	"log"
	"os"

	"github.com/aleksandradidenko/go_final_project/pkg/db"
	"github.com/aleksandradidenko/go_final_project/pkg/server"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatal(err)
	}

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}
