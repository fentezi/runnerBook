package main

import (
	"github.com/fentezi/runnerBook/config"
	"github.com/fentezi/runnerBook/server"
	"log"
)

func main() {
	log.Println("Starting Runners App")
	log.Println("Initializing configuration")
	config := config.InitConfig("runners")
	log.Println("Initializing database")
	dbHandler := server.InitDatabase(config)
	log.Println("Initializing HTTP server")
	httpServer := server.InitHttpServer(config, dbHandler)
	httpServer.Start()
}
