
## Actor Objects 

> Objects with an active runloop / channels comunication.
> An instances is like a thread object.


- Server 

    rapresents server instance

- PubSubEngine (Router)
    
    responsable of publish/subscribe logic, we can have different implementations: 
    ram, redis, kafka, google pubsub, ...
    
        Server >>use>> PubSubEngine
 
- Client 
    
    one instance per client session

## Networking and protocol handling

- ConnectionListener 

    tcp or websocket listener
    expose "onNewConnection" Event 
        with Connection object prepared

- Connection 
    
    A wrapper over a tcp or websocket connection.
    
        (Client#1)-----[Connection#1]----->(Server)<------[Connection#2]-----(Client#2)
    
        Connecton >>use>> MqttProtocolHandler >>use>> TcpConn        
    
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
    

## Main Flow

- new Server ( -> new PubSubEngine )
- Listen to TCP connection
    - On Connect
        - new Client Object
        - new Connection ( -> new MqttProtocolHandler ) 
            with Client and Server instances
        - Start Connection read/write tcp socket loop
        
            connection interacts with client and server throught channels
        
            connection read from tcp and notify client/server 
            and than write to tcp according to mqtt protocol
            
        