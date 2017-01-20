# Mqtt Broker POC

## Stats

    publishchan := make(chan *packets.PublishPacket)
    addchan := make(chan Client)
    rmchan := make(chan Client)
    ------------------------------------------------------------------
    Test clients: 10, num msg: 10000, qos: 0 sleep: 4000 millis.
    ------------------------------------------------------------------
    Time elapsed: 8391 millis.
    Messages sent: 100000 messages.
    Messages received: 225307 messages.
    Throughput Published: 12500.0 msg/sec.
    Throughput Received: 28163.625 msg/sec.
    Throughput Total: 40663.625 msg/sec.
    ------------------------------------------------------------------
    
    
    
    
    publishchan := make(chan *packets.PublishPacket, 100)
    addchan := make(chan Client, 10)
    rmchan := make(chan Client, 10)
    ------------------------------------------------------------------
    Test clients: 10, num msg: 10000, qos: 0 sleep: 4000 millis.
    ------------------------------------------------------------------
    Time elapsed: 8416 millis.
    Messages sent: 100000 messages.
    Messages received: 253274 messages.
    Throughput Published: 12500.0 msg/sec.
    Throughput Received: 31659.625 msg/sec.
    Throughput Total: 44159.625 msg/sec.
    ------------------------------------------------------------------
    
    
    publishchan := make(chan *packets.PublishPacket, 100000)
    addchan := make(chan Client, 1000)
    rmchan := make(chan Client, 1000)
    ------------------------------------------------------------------
    Test clients: 10, num msg: 10000, qos: 0 sleep: 4000 millis.
    ------------------------------------------------------------------
    Time elapsed: 8341 millis.
    Messages sent: 100000 messages.
    Messages received: 217184 messages.
    Throughput Published: 12500.0 msg/sec.
    Throughput Received: 27148.0 msg/sec.
    Throughput Total: 39648.0 msg/sec.
    ------------------------------------------------------------------


	publishchan := make(chan *packets.PublishPacket, 10)
	addchan := make(chan Client, 5)
	rmchan := make(chan Client, 5)
    ------------------------------------------------------------------
    Test clients: 10, num msg: 10000, qos: 0 sleep: 4000 millis.
    ------------------------------------------------------------------
    Time elapsed: 8238 millis.
    Messages sent: 100000 messages.
    Messages received: 296310 messages.
    Throughput Published: 12500.0 msg/sec.
    Throughput Received: 37039.0 msg/sec.
    Throughput Total: 49539.0 msg/sec.
    ------------------------------------------------------------------
    
    
    publishchan := make(chan *packets.PublishPacket)
    addchan := make(chan Client)
    rmchan := make(chan Client)
    ------------------------------------------------------------------
    Test clients: 10, num msg: 10000, qos: 0 sleep: 4000 millis.
    ------------------------------------------------------------------
    Time elapsed: 8144 millis.
    Messages sent: 100000 messages.
    Messages received: 289701 messages.
    Throughput Published: 12500.0 msg/sec.
    Throughput Received: 36213.125 msg/sec.
    Throughput Total: 48713.125 msg/sec.
    ------------------------------------------------------------------
    
    
## After optimization of regexp calculation

    ------------------------------------------------------------------
    Test clients: 10, num msg: 10000, qos: 0 sleep: 3000 millis.
    ------------------------------------------------------------------
    Time elapsed: 8200 millis.
    Messages sent: 100000 messages.
    Messages received: 987690 messages.
    Throughput Published: 12500.0 msg/sec.
    Throughput Received: 123464.0 msg/sec.
    Throughput Total: 135964.0 msg/sec.
    ------------------------------------------------------------------