package main

import (
	"context"

	"github.com/brendenehlers/todo-microservice/http"
	"github.com/brendenehlers/todo-microservice/memory"
	"github.com/brendenehlers/todo-microservice/slogger"
)

func main() {
	log := slogger.New()
	repo := memory.New(log)

	ctx := context.Background()
	server, _ := http.CreateHTTPServer(&http.HTTPServerConfig{
		Addr: ":8080",
		Repo: repo,
		Ctx:  ctx,
		Log:  log,
	})

	server.Run()
}
