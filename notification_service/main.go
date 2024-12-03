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

	notificationApp, err := NewNotificationApp(stopServerChannel)

	if err != nil {
		fmt.Println("Failed to start costumers service")
		return
	}

	defer notificationApp.Close()

	notificationApp.StartWorkers()

	<-stopServerChannel
}
