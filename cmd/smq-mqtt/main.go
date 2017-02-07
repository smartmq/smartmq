package main

import (
	"flag"
	"fmt"
	"github.com/smartmq/smartmq"
	"github.com/smartmq/smartmq/cmd/smq-mqtt/mqtt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var _redisURL string

func main() {

	url := flag.String("redis-url", "redis://127.0.0.1:6379", "redis url")
	redisURL, exists := os.LookupEnv("REDIS_URL")
	if !exists {
		redisURL = *url
	}
	_redisURL = redisURL

	mq := smartmq.New(_redisURL, false)

	router := mqtt.NewRouter(mq)
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
		mq.Close()
		os.Exit(0)
	}
}
