package main

import (
	"flag"
	"fmt"
	"github.com/pborman/uuid"
	"github.com/smartmq/smartmq"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	url := flag.String("url", "redis://127.0.0.1:6379", "redis url")
	chn := flag.String("c", "smartmq", "channel (or topic)")
	sid := flag.String("s", uuid.New(), "subscription queue id")
	enableLogging := flag.Bool("l", false, "activate trace logs")
	flag.Parse()

	redisURL, exists := os.LookupEnv("REDIS_URL")
	if !exists {
		redisURL = *url
	}
	channel := *chn
	subid := *sid

	mq := smartmq.New(redisURL, *enableLogging)

	subs := mq.Channel(channel).AddSubscription(subid).StartConsuming(func(key string, val string) {
		fmt.Printf("%v << %v\n", key, val)
	})

	// capture ctrl+c
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	select {
	case <-c:
		if *enableLogging {
			log.Println("Shutting down ...")
		}
		//subs.StopConsuming()
		//subs.RemoveSubscription()
		subs.Close()
		if *enableLogging {
			log.Println("Unsubscribed")
		}
		mq.Close()
		if *enableLogging {
			log.Println("Connection closed")
		}
		os.Exit(0)
	}
}
