package mqtt

import (
	"fmt"
	"github.com/smartmq/smartmq/cmd/smq-mqtt/packets"
	"golang.org/x/net/websocket"
	"log"
	"net"
	"net/http"
	"time"
)

func StartTcpServer(laddr string, router *Router) {
	ln, err := net.Listen("tcp", laddr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Listen on %s\n", laddr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn, router)
	}
}

func StartWebsocketServer(addr string, router *Router) {
	fmt.Printf("Listen on %s (websocket)\n", addr)
	acceptConnection := func(ws *websocket.Conn) {
		handleConnection(ws, router)
	}
	handler := websocket.Handler(acceptConnection)
	http.Handle("/mqtt", handler)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handleConnection(conn net.Conn, router *Router) {
	client := NewClient(conn)

	defer func() {
		log.Printf("Closing connection from %v (clientID: %s)...\n", conn.RemoteAddr(), client.ID)
		conn.Close()
		if client.CleanSession {
			router.Unsubscribe(client)
			router.Disconnect(client)
		}
		log.Printf("Connection from %v closed (clientID: %s).\n", conn.RemoteAddr(), client.ID)
	}()

	log.Printf("Connection from %v.\n", conn.RemoteAddr())

	go handleMqttProtocol(router, client)
	client.Start(router)
}

func handleMqttProtocol(router *Router, client *Client) {
	for {
		select {
		case cp := <-client.incoming:
			//log.Printf(">> %v \n", cp.String())

			//msgType := cp.GetMessageType()
			switch cp.(type) {

			case *packets.ConnectPacket:
				connectMsg := cp.(*packets.ConnectPacket)
				client.ID = connectMsg.ClientIdentifier
				client.CleanSession = connectMsg.CleanSession
				client.Keepalive = NewKeepalive(connectMsg.Keepalive)
				client.Keepalive.ExpiredCallback = func(t time.Time) {
					log.Printf("Keepalive time exausted for client: %s", client.ID)
					disconnectAbnormally(client, router)
				}

				if connectMsg.WillFlag {
					will := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
					will.TopicName = connectMsg.WillTopic
					will.Payload = connectMsg.WillMessage
					will.Qos = connectMsg.WillQos
					client.WillMessage = will
					if client.WillMessage != nil {
						log.Printf("Client %s has will message\n", client.ID)
					}
				}
				proto := connectMsg.ProtocolName
				if proto == "MQTT" {
					//log.Printf("MQTT 3.1.1 (%s)\n", proto)
				} else if proto == "MQIsdp" {
					//log.Printf("MQTT 3.1 (%s)\n", proto)
				} else {
					log.Printf("Wrong protocol (%s)\n", proto)
					disconnectAbnormally(client, router)
					break
				}
				connackMsg := packets.NewControlPacket(packets.Connack).(*packets.ConnackPacket)

				if router.Connected(client) {
					log.Printf("Client alredy connected %s", client.ID)
					connackMsg.SessionPresent = true

					// recupera la sessione trovata ...
					oldClient := router.GetConnected(client.ID)
					sameConnection := oldClient.Conn == client.Conn
					//oldClientConnected := oldClient.IsConnected()

					//log.Printf("Old Client connected: %s, same connection: %s", oldClientConnected, sameConnection)
					log.Printf("Old Client same connection: %s", sameConnection)

					if sameConnection {
						disconnectAbnormally(client, router)
						break
					} else {
						if !connectMsg.CleanSession {
							oldClient.CopyTo(client)
							// swap clients ...
							router.Disconnect(oldClient)
						}
						router.Connect(client)
					}
				} else {
					router.Connect(client)
					log.Printf("New Client connected %s", client.ID)
					connackMsg.SessionPresent = false
				}

				client.outgoing <- connackMsg

				// start keepalive timer for this client
				client.Keepalive.Start()

				client.FlushQueuedMessages()

				break
			case *packets.SubscribePacket:
				client.Keepalive.Reset()

				subackMsg := packets.NewControlPacket(packets.Suback).(*packets.SubackPacket)
				subackMsg.MessageID = cp.Details().MessageID
				router.Subscribe(client)

				subMsg := cp.(*packets.SubscribePacket)
				topics := subMsg.Topics
				qoss := subMsg.Qoss

				for i := 0; i < len(topics); i++ {
					topic := topics[i]
					var qos byte
					if i < len(qoss) {
						qos = qoss[i]
					} else {
						qos = 0x0
					}
					s := NewSubscription(topic, qos)
					client.Subscriptions[topic] = s
				}

				client.outgoing <- subackMsg

				router.RepublishRetainedMessages(client, subMsg)
				break
			case *packets.UnsubscribePacket:
				client.Keepalive.Reset()
				router.Unsubscribe(client)
				unsubackMsg := packets.NewControlPacket(packets.Unsuback).(*packets.UnsubackPacket)
				unsubackMsg.MessageID = cp.Details().MessageID
				client.outgoing <- unsubackMsg
				break
			case *packets.DisconnectPacket:
				client.Keepalive.Reset()
				disconnect(client, router)
				break
			case *packets.PublishPacket:

				client.Keepalive.Reset()
				qos := cp.Details().Qos
				pubMsg := cp.(*packets.PublishPacket)
				//log.Printf("%v: PUBLISH qos(%v)", client.ID, pubMsg.Qos)

				switch qos {
				case 0:
					router.Publish(pubMsg)
					break
				case 1:
					router.Publish(pubMsg)
					pubAck := packets.NewControlPacket(packets.Puback).(*packets.PubackPacket)
					pubAck.MessageID = cp.Details().MessageID
					client.outgoing <- pubAck
					break
				case 2:
					router.Publish(pubMsg)
					pubRec := packets.NewControlPacket(packets.Pubrec).(*packets.PubrecPacket)
					pubRec.MessageID = cp.Details().MessageID
					client.outgoing <- pubRec
					break
				}

			case *packets.PubackPacket:
				client.Keepalive.Reset()
				// PUBACK (receiver) response to PUBLISH in QOS=1
				client.queue.DequeueMessage()

				//pubAck := cp.(*packets.PubackPacket)
				//log.Printf("%v: PUBACK qos(%v)", client.ID, pubAck.Qos)
				break

			case *packets.PubrecPacket:
				client.Keepalive.Reset()
				// PUBREC (receiver) response to PUBLISH in QOS=2
				//pubRec := cp.(*packets.PubrecPacket)
				//log.Printf("%v: PUBREC qos(%v)", client.ID, pubRec.Qos)
				pubRel := packets.NewControlPacket(packets.Pubrel).(*packets.PubrelPacket)
				pubRel.MessageID = cp.Details().MessageID
				client.outgoing <- pubRel
				break

			case *packets.PubrelPacket:
				client.Keepalive.Reset()
				// PUBREL (sender) response to PUBREC in QOS=2
				//pubRel := cp.(*packets.PubrelPacket)
				//log.Printf("%v: PUBREL qos(%v)", client.ID, pubRel.Qos)
				pubComp := packets.NewControlPacket(packets.Pubcomp).(*packets.PubcompPacket)
				pubComp.MessageID = cp.Details().MessageID
				client.outgoing <- pubComp
				break

			case *packets.PubcompPacket:
				client.Keepalive.Reset()
				// PUBCOMP (receiver) response to PUBREL in QOS=2
				client.queue.DequeueMessage()

				//pubComp := cp.(*packets.PubcompPacket)
				//log.Printf("%v: PUBCOMP qos(%v)", client.ID, pubComp.Qos)
				break

			case *packets.PingreqPacket:
				client.Keepalive.Reset()
				pingresp := packets.NewControlPacket(packets.Pingresp).(*packets.PingrespPacket)
				client.outgoing <- pingresp
				break

			default:
				disconnectAbnormally(client, router)
			}

		}
	}
}

func disconnect(client *Client, router *Router) {
	if client.CleanSession {
		router.Unsubscribe(client)
		router.Disconnect(client)
	}
	client.Quit()
}

func disconnectAbnormally(client *Client, router *Router) {
	log.Printf("Disconnect abnormally from client: %s\n", client.ID)
	handleWillMessage(client, router)
	router.Unsubscribe(client)
	router.Disconnect(client)
	client.Quit()
}

func handleWillMessage(client *Client, router *Router) {
	log.Printf("publish will message of client %s if present ...\n", client.ID)
	// publish will message if present ...
	if client.WillMessage != nil {
		log.Printf("will message topic[%s]\n", client.WillMessage.TopicName)
		client.WillMessage.MessageID = 1
		router.Publish(client.WillMessage)
	}
}
