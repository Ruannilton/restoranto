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

	verificationService, err := NewContactVerificationApp(stopServerChannel)

	if err != nil {
		fmt.Println("Failed to start costumers service")
		return
	}

	defer verificationService.Close()

	verificationService.StartWorkers()

	<-stopServerChannel
}
