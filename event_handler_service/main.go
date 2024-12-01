package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	stopServerChannel := make(chan os.Signal, 1)
	signal.Notify(stopServerChannel, syscall.SIGINT, syscall.SIGTERM)

	if err != nil {
		fmt.Println("Failed to read env vars")
		return
	}

	eventHandlerApp, err := NewEventHandlerApp(stopServerChannel)

	if err != nil {
		fmt.Println("Failed to start costumers service")
		return
	}

	defer eventHandlerApp.Close()

	eventHandlerApp.StartWorkers()

	<-stopServerChannel
}
