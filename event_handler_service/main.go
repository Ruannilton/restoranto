package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	stopServerChannel := make(chan os.Signal, 1)
	signal.Notify(stopServerChannel, syscall.SIGINT, syscall.SIGTERM)

	eventHandlerApp, err := NewEventHandlerApp(stopServerChannel)

	if err != nil {
		fmt.Println("Failed to start costumers service")
		return
	}

	defer eventHandlerApp.Close()

	eventHandlerApp.StartWorkers()

	<-stopServerChannel
}
