package main

import (
	"os"
	"github.com/smartmq/smartmq"
	"flag"
)

func main() {

	url := flag.String("url", "redis://127.0.0.1:6379", "redis url")
	chn := flag.String("c", "", "channel")
	msg := flag.String("m", "", "message to publish")
	enableLogging := flag.Bool("l", false, "activate trace logs")
	flag.Parse()


	redisURL, exists := os.LookupEnv("REDIS_URL")
	if !exists {
		redisURL = *url
	}
	channel := *chn
	message := *msg

	mq := smartmq.New(redisURL, *enableLogging)
	defer mq.Close()

	mq.Channel(channel).Publish([]byte(message))
}
