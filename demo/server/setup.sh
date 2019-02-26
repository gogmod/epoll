#!/bin/bash
docker run -v $(pwd)/server:/go/server -p 6060:6060 --name tcp_server -d golang /go/server
