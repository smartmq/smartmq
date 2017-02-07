package smartmq

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

type SmartMQ struct {
	redisURL      string
	c             redis.Conn
	enableLogging bool
}

type Channel struct {
	mq      *SmartMQ
	channel string
	key     string
}

type Subscription struct {
	channel       *Channel
	subscription  string
	key           string
	stopComsuming chan bool
}

type OnMessageFn func(key string, val string)

func New(redisURL string, enableLogging bool) *SmartMQ {
	mq := &SmartMQ{
		redisURL:      redisURL,
		enableLogging: enableLogging,
	}
	return mq.Open()
}

func (mq *SmartMQ) newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(mq.redisURL)
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func (mq *SmartMQ) Open() *SmartMQ {
	mq.c = mq.newPool().Get()
	if mq.enableLogging {
		logger := log.New(os.Stderr, "SmartMQ", log.LstdFlags)
		mq.c = redis.NewLoggingConn(mq.c, logger, mq.redisURL)
	}
	return mq
}

func (mq *SmartMQ) Close() {
	err := mq.c.Close()
	handleError(err)
}

func (mq *SmartMQ) Channel(channel string) *Channel {
	return &Channel{
		mq:      mq,
		channel: channel,
		key:     key(channel, "subscriptions"),
	}
}

func (ch *Channel) MQ() *SmartMQ {
	return ch.mq
}

//func (ch *Channel) _publishWithRef(message string) {
//	c := ch.mq.c
//	channel := ch.channel
//	val, err := redis.Int64(c.Do("INCR", key(channel, "nextid")))
//	var messageID string
//	messageID = strconv.FormatInt(val, 10)
//	err = c.Send("SET", key(channel, "messages", messageID), message, "EX", 10)
//	handleError(err)
//	err = c.Flush()
//	handleError(err)
//
//	subscriptions, err := redis.Strings(c.Do("SMEMBERS", ch.key))
//
//	for i := range subscriptions {
//		sub := subscriptions[i]
//		key := key(channel, sub, "messages")
//		debug("PUB", key)
//		err = c.Send("LPUSH", key, messageID)
//		handleError(err)
//	}
//	err = c.Flush()
//	handleError(err)
//}

func (ch *Channel) Publish(message []byte) {
	var err error
	c := ch.mq.c
	k := ch.key
	subscriptions, err := redis.Strings(c.Do("SMEMBERS", k))
	for i := range subscriptions {
		sub := subscriptions[i]
		key := key(ch.channel, sub, "messages")
		c.Send("LPUSH", key, message)
		//c.Send("LTRIM", key, 0, 9)
	}
	c.Send("LPUSH", key, message)
	err = c.Flush()
	handleError(err)
}

func (ch *Channel) Subscriptions() []*Subscription {
	c := ch.mq.c
	k := ch.key
	subscriptions, err := redis.Strings(c.Do("SMEMBERS", k))
	handleError(err)
	list := make([]*Subscription, len(subscriptions))
	for i := range subscriptions {
		sub := subscriptions[i]
		subs := &Subscription{
			channel:       ch,
			subscription:  sub,
			key:           key(ch.channel, sub, "messages"),
			stopComsuming: make(chan bool),
		}
		list[i] = subs
	}
	return list
}

func (ch *Channel) AddSubscription(subscription string) *Subscription {
	c := ch.mq.c
	err := c.Send("SADD", ch.key, subscription)
	err = c.Flush()
	if err != nil {
		log.Fatal(err)
	}
	return ch.Subscription(subscription)
}
func (ch *Channel) Subscription(subscription string) *Subscription {
	sub := &Subscription{
		channel:       ch,
		subscription:  subscription,
		key:           key(ch.channel, subscription, "messages"),
		stopComsuming: make(chan bool),
	}
	return sub
}

func (subs *Subscription) StartConsuming(fn OnMessageFn) *Subscription {
	c := subs.channel.mq.c

	go func() {
		k := subs.key
		for {
			select {
			case quit := <-subs.stopComsuming:
				if quit {
					//log.Println("stop listening of new messages ...")
					return
				}
			default:
				val, _ := redis.Strings(c.Do("BRPOP", k, 1))
				if len(val) == 2 {
					fn(val[0], val[1])
				}
			}

		}
	}()
	return subs
}

func (subs *Subscription) StopConsuming() *Subscription {
	subs.stopComsuming <- true
	return subs
}

func (subs *Subscription) GetMessage() []byte {
	c := subs.channel.mq.c
	k := subs.key

	//val, _ := redis.Strings(c.Do("BRPOP", k, 10))
	//if len(val) == 2 {
	//	return val[1]
	//}
	val, _ := redis.ByteSlices(c.Do("RPOP", k))
	if len(val) == 2 {
		return val[1]
	}
	return []byte{}
}
func (subs *Subscription) GetMessageBlocking() []byte {
	c := subs.channel.mq.c
	k := subs.key
	val, _ := redis.ByteSlices(c.Do("BRPOP", k, 10))
	if len(val) == 2 {
		return val[1]
	}
	return []byte{}
}

func (subs *Subscription) Close() *Channel {
	subs.StopConsuming()
	subs.RemoveSubscription()
	return subs.channel
}

func (subs *Subscription) RemoveSubscription() {
	c := subs.channel.mq.c
	k := subs.channel.key
	_, err := c.Do("SREM", k, subs.subscription)
	err = c.Flush()
	handleError(err)
}

func (subs *Subscription) PurgeSubscriptionQueue() {
	c := subs.channel.mq.c
	subscriptionKey := subs.key
	_, err := c.Do("DEL", subscriptionKey)
	err = c.Flush()
	handleError(err)
}

func (subs *Subscription) Channel() *Channel {
	return subs.channel
}

func (subs *Subscription) ToString() string {
	return fmt.Sprintf("%s (%s)", subs.subscription, subs.key)
}

func key(s ...string) string {
	return strings.Join(s, ".")
}

func debug(id string, val interface{}) {
	log.Printf("DEBUG %s  %v (%v)\n", id, val, reflect.TypeOf(val))
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
