#!/bin/bash
    docker run -v $(pwd)/client:/go/server --name tcp_server -d golang /go/server
