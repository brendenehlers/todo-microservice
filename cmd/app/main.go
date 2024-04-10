package main

import (
	"context"

	"github.com/brendenehlers/todo-microservice/http"
	"github.com/brendenehlers/todo-microservice/logger"
	"github.com/brendenehlers/todo-microservice/memory"
)

func main() {
	log := &logger.Logger{}
	repo := memory.New()

	server, _ := http.CreateHTTPServer(&http.HTTPServerConfig{
		Addr: ":8080",
		Repo: repo,
		Ctx:  context.Background(),
		Log:  log,
	})

	server.Run()
}
