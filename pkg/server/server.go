package server

import (
	"net/http"
	"os"

	"github.com/aleksandradidenko/go_final_project/pkg/api"
)

func Run() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	api.Init()

	http.Handle("/", http.FileServer(http.Dir("web")))

	return http.ListenAndServe(":"+port, nil)
}
