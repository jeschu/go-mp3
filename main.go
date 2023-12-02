package main

import (
	"go-mp3/engine"
	"go-mp3/engine/controller"
	"go-mp3/env"
	"go-mp3/lang"
	"log"
)

func main() {
	config := env.InitFromEnv()

	server, err := engine.NewServer().Addr(config.Addr).ManagementAddr(config.ManagementAddr).WithRequestId(config.WithRequestId).Build()
	if err != nil {
		log.Fatal(err)
	}
	server.RegisterController(controller.VersionController{})
	server.Start()
	server.RegisterController(controller.NewLibraryController())
	lang.WaitForIntOrTerm()
	server.Shutdown()
}
