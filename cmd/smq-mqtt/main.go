package main

import (
	"fmt"
	"github.com/smartmq/smartmq/cmd/smq-mqtt/mqtt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	router := mqtt.NewRouter()
	go router.Start()

	go mqtt.StartTcpServer(":1883", router)
	go mqtt.StartWebsocketServer(":11883", router)

	go http.ListenAndServe("localhost:6060", nil)

	// capture ctrl+c
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	select {
	case <-c:
		fmt.Println("Shutting down ...")
		os.Exit(0)
	}
}
