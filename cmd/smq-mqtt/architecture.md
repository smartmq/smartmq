
## Actor Objects 

> Objects with an active runloop / channels comunication.
> An instances is like a thread object.

- Server 

    probabily the main instance who start all others 
 
- Client 
    
    one instance per client session

- PubSubEngine (Router)
    
    responsable of publish/subscribe logic, we can have different implementations: 
    ram, redis, kafka, google pubsub, ...

- Connection 
    
    A wrapper over a tcp or websocket connection.
    
- MqttProtocolHandler

    Encode and Decode MQTT protocol, and let communicate client with server
    all throught chan golang api, without expose mqtt structures or socket api. 
    
    Those are methods to be refactored into this object:
    
        - client.go: Client.Read() 
        - client.go: Client.Write()
        - server.go: handleMqttProtocol(router *Router, client *Client)



## Semantic data Objects (no active thread)

> Simple Classes with private data and public methods to manipulate data.

- Topic 
    
    Topic string wrapper
    
- Payload 

    Payload byte[] wrapper

- Message 

    composite of topic + payload
    
- Queue 

    we need queue with fifo capabilities, so implements a solid one
        
- Subscription

    Wrap a client subscription, with topic filter and qos info.
    This is part of client session or state
    
