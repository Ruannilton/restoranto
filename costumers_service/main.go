package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	stopServerChannel := make(chan os.Signal, 1)
	signal.Notify(stopServerChannel, syscall.SIGINT, syscall.SIGTERM)

	costumerService, err := NewCostumersServiceApp(stopServerChannel)

	if err != nil {
		fmt.Println("Failed to start costumers service")
		return
	}

	router := mux.NewRouter()

	costumerService.RegisterControllers(router)
	costumerService.StartWorkers()

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
