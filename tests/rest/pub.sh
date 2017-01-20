#!/bin/bash

for i in `seq 1 5000`;
do
    curl -XPUT -d "una prova $(date) ${i}
" http://localhost:8080/publish?channel=chan/mqtt/topic
done

#curl -XPUT -d '{"chan":"chan/mqtt/topic","msg":"ciao mondo!"}' http://localhost:8080/publish?channel=chan/mqtt/topic
#curl -XPUT -d '@email_fili_logo.gif' http://localhost:8080/publish?channel=chan/mqtt/topic
