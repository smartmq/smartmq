package mqtt

import (
	"log"
	"testing"
	"time"
)

func TestNewKeepalive(t *testing.T) {
	k := NewKeepalive(1)

	k.ExpiredCallback = func(t time.Time) {
		log.Printf("Keepalived expired at: %v\n", t)
		k.Stop()
	}

	log.Println("Start")
	k.Start()

	time.Sleep(2 * time.Second)

	log.Println("Reset 1")
	k.Reset()

	time.Sleep(2 * time.Second)

	log.Println("Reset 2")
	k.Reset()

	time.Sleep(5 * time.Second)

	// dovrebbe scattare qui

	k.Reset()
	log.Println("Reset 3")

	time.Sleep(5 * time.Second)

	k.Stop()
	log.Println("Stop")
}
