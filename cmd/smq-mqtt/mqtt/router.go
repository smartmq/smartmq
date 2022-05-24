package mqtt

import (
	"github.com/smartmq/smartmq/cmd/smq-mqtt/packets"
	"sync"
)

type Router struct {
	clients          map[string]*Client
	connectedClients map[string]*Client

	publishChan     chan *packets.PublishPacket
	subscribeChan   chan *Client
	unsubscribeChan chan *Client

	retainStore map[string]*packets.PublishPacket
	queues      map[string]*Queue

	lock sync.RWMutex
}

func NewRouter() *Router {
	return &Router{
		clients:          make(map[string]*Client),
		connectedClients: make(map[string]*Client),
		publishChan:      make(chan *packets.PublishPacket),
		subscribeChan:    make(chan *Client),
		unsubscribeChan:  make(chan *Client),
		retainStore:      make(map[string]*packets.PublishPacket),
		queues:           make(map[string]*Queue),
		lock:             sync.RWMutex{},
	}
}

func (router *Router) addClient(c *Client) {
	router.lock.Lock()
	router.clients[c.ID] = c
	router.lock.Unlock()
}
func (router *Router) removeClient(c *Client) {
	router.lock.Lock()
	delete(router.clients, c.ID)
	router.lock.Unlock()
}

func (router *Router) Subscribe(c *Client) {
	router.subscribeChan <- c
}
func (router *Router) Unsubscribe(c *Client) {
	router.unsubscribeChan <- c
}
func (router *Router) Publish(pubMsg *packets.PublishPacket) {
	router.publishChan <- pubMsg
}

func (router *Router) RepublishRetainedMessages(c *Client, subMsg *packets.SubscribePacket) {
	msgid := subMsg.MessageID
	for _, pubMsg := range router.retainStore {
		msgid++
		msg := pubMsg.Copy()
		msg.FixedHeader = pubMsg.FixedHeader
		msg.MessageID = msgid
		sendMessageToClientIfMatch(c, msg)
	}
}

//func (router *Router) Start() {
//	go router.handleEvents()
//}

func (router *Router) Start() {
	clients := router.clients
	for {
		select {
		case msg := <-router.publishChan:
			// save if retain
			if msg.Retain {
				payloadLen := len(msg.Payload)
				//log.Printf("payload retained: %s\n", payloadLen)
				if payloadLen == 0 {
					//log.Printf("payload retained empty: %s, delete retained message\n", payloadLen)
					delete(router.retainStore, msg.TopicName)
				} else {
					router.retainStore[msg.TopicName] = msg
				}
			}
			// foreach client send message that match topic
			for _, c := range clients {
				sendMessageToClientIfMatch(c, msg)
			}
		case c := <-router.subscribeChan:
			router.addClient(c)
		case c := <-router.unsubscribeChan:
			router.removeClient(c)
		}
	}
}

func sendMessageToClientIfMatch(client *Client, msg *packets.PublishPacket) {
	// check subscription match
	publishingTopic := msg.TopicName
	subscribed, subscriptionQos := client.IsSubscribed(publishingTopic)
	if subscribed {
		//log.Printf("sending message to client: %s topic: %s (%s)\n", client.ID, publishingTopic, subscribed)
		// TODO: look at mqtt specs for qos handling...
		msg.Qos = subscriptionQos
		if msg.MessageID <= 0 {
			//log.Fatal("MessageID is 0 !")
			msg.MessageID = 1
		}
		if msg.Qos > 0 {
			client.queue.EnqueueMessage(msg)
		}
		client.WritePublishMessage(msg)
	}
}

func (r *Router) Connect(c *Client) {
	r.lock.Lock()
	r.connectedClients[c.ID] = c
	r.lock.Unlock()
}
func (r *Router) Disconnect(c *Client) {
	r.lock.Lock()
	delete(r.connectedClients, c.ID)
	r.lock.Unlock()
}
func (r *Router) Connected(c *Client) bool {
	r.lock.RLock()
	_, present := r.connectedClients[c.ID]
	r.lock.RUnlock()
	return present
}
func (r *Router) GetConnected(client_id string) *Client {
	r.lock.RLock()
	c, _ := r.connectedClients[client_id]
	r.lock.RUnlock()
	return c
}
