# smart-mq 
_A Pub/Sub wrapper service over redis_

## Concepts

### Channel
A channel is a stream of data.
It is a simple string.
A publisher must always specify a channel to publish stream of data.

### Subscription
A subscription is an abstract subscription to a channel.
It is a simple string.
In server, a subscription is essentually a queue, who receive data from channel.
A subscription is always attached to only one channel.
A subscription can exists without a consumer, 
it receive data for consumenrs that are offline.
If 2 or more client subscribe to same subscription, 
the messages are dispatched one per client, so client receive partial data.
If a client want receive all data from a channel, 
it must create an unique subscription.

## Verbs

### Subscription phase
- register a new subscription to a channel
- consume data from subscription of a channel
- remove subscription from channel
- purge a subscription of a channel

### Publish phase
- publish a message into a channel


## REST Protocol

    POST   /mq/:channel               (publish new message)
    POST   /mq/:channel/:subscription (add new subscription)
    GET    /mq/:channel/:subscription (consume messages from subscription)
    DELETE /mq/:channel/:subscription (remove a subscription)

## MQTT Protocol

## gRPC Protocol

