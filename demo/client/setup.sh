#!/bin/bash
CONNECTIONS=$1
REPLICAS=$2
IP=$3
#go build --tags "static netgo" -o client client.go
for (( c=0; c<${REPLICAS}; c++ ))
do
    docker run -v $(pwd)/client:/go/client --name 1mclient_$c -d golang /go/client -conn=${CONNECTIONS} -ip=${IP}
done
