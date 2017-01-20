package mqtt

import (
	"time"
)

type Keepalive struct {
	Duration        time.Duration
	ExpiredCallback func(t time.Time)
	expired         bool
	ticker          *time.Ticker
}

func NewKeepalive(seconds uint16) *Keepalive {
	if seconds == 0 {
		return &Keepalive{
			Duration: 0,
			expired:  false,
		}
	} else {
		duration := time.Millisecond * time.Duration(seconds*1500)
		return &Keepalive{
			Duration: duration,
			expired:  false,
		}
	}
}

func (k *Keepalive) Reset() {
	k.expired = false
}
func (k *Keepalive) Start() {
	/*
	 * If the Keep Alive value is non-zero and the Server does not receive a Control Packet from the Client
	 * within one and a half times the Keep Alive time period, it MUST disconnect
	 */
	//t := time.Now()
	//log.Printf("Keepalive start: %v\n", t)
	if k.Duration == 0 {
		return
	}
	k.expired = true
	k.ticker = time.NewTicker(k.Duration)
	go func() {
		for t := range k.ticker.C {
			//log.Printf("Keepalive tiker: %v\n", t)
			if k.expired {
				k.ExpiredCallback(t)
			} else {
				k.expired = true
			}
		}
	}()
}

func (k *Keepalive) Stop() {
	if k.ticker != nil {
		k.ticker.Stop()
	}
}
