#!/bin/bash
# add a subscription named "sub"
#curl -XPOST http://localhost:8080/subscribe/sub1?channel=chan/mqtt/topic

#sleep 5

# get messages
while :
do
    curl -XGET http://localhost:8080/subscribe/sub1?channel=chan/mqtt/topic
done
