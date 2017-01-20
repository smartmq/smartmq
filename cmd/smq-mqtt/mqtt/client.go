package mqtt

import (
	"bufio"
	"github.com/smartmq/smartmq/cmd/smq-mqtt/packets"
	"log"
	"net"
)

type Client struct {
	Conn          net.Conn
	ID            string
	CleanSession  bool
	WillMessage   *packets.PublishPacket
	Keepalive     *Keepalive
	reader        *bufio.Reader
	writer        *bufio.Writer
	queue         *Queue
	incoming      chan packets.ControlPacket
	outgoing      chan packets.ControlPacket
	Subscriptions map[string]*Subscription
	quit          chan bool
}

func NewClient(connection net.Conn) *Client {
	writer := bufio.NewWriter(connection)
	//writer := bufio.NewWriterSize(connection, (4096*16))
	reader := bufio.NewReader(connection)

	client := &Client{
		Conn:          connection,
		reader:        reader,
		writer:        writer,
		queue:         NewQueue(),
		incoming:      make(chan packets.ControlPacket),
		outgoing:      make(chan packets.ControlPacket),
		Subscriptions: make(map[string]*Subscription),
		quit:          make(chan bool),
	}
	return client
}

func (client *Client) Read() {
	for {
		//log.Println("reading ...")
		data, err := packets.ReadPacket(client.reader)
		if err != nil {
			break
		}
		client.incoming <- data
	}
}
func (client *Client) Write() {
	for data := range client.outgoing {
		//log.Println("writing ...")
		//log.Printf(">> %v \n", data)
		err := data.Write(client.writer)
		if err != nil {
			//log.Println("write ", client.ID, " ",  err)
		}
		err2 := client.writer.Flush()
		if err2 != nil {
			//log.Println("flush ", client.ID, " ", err)

			// enqueue messsage if error and cleanSession is false
			if !client.CleanSession {
				switch data.(type) {
				case *packets.PublishPacket:
					pubMsg := data.(*packets.PublishPacket)
					if pubMsg.Qos == 1 || pubMsg.Qos == 2 {
						client.queue.EnqueueMessage(pubMsg)
					}
					break
				}
			}
		}
	}
}
func (client *Client) Quit() {
	log.Printf("Quitting client %s ...", client.ID)
	client.Keepalive.Stop()
	client.Conn.Close()
	client.quit <- true
}
func (client *Client) waitForQuit() {
	select {
	case ret := <-client.quit:
		if ret {
			return
		}
	}
}
func (c *Client) IsSubscribed(publishingTopic string) (bool, byte) {
	var ret bool
	var qos byte
	ret = false
	qos = 0x00
	for _, subscription := range c.Subscriptions {
		if subscription.IsSubscribed(publishingTopic) {
			ret = true
			if qos < subscription.Qos {
				qos = subscription.Qos
			}
		}
	}
	return ret, qos
}
func (c *Client) WritePublishMessage(msg *packets.PublishPacket) {
	c.outgoing <- msg
}

func (client *Client) Start(router *Router) {
	go client.Read()
	go client.Write()
	client.waitForQuit()
}

func (client *Client) CopyTo(other *Client) {
	other.Subscriptions = client.Subscriptions
	other.queue = client.queue
}

func (client *Client) FlushQueuedMessages() {
	queueSize := client.queue.Size()
	if queueSize > 0 {
		log.Printf("Flush queued messages (%v)", queueSize)
		for msg := client.queue.DequeueMessage(); msg != nil; msg = client.queue.DequeueMessage() {
			client.WritePublishMessage(msg.(*packets.PublishPacket))
		}
	}
}
