package main

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"

	"github.com/gogmod/epoll"
)

var epoller *epoll.Epoll

//docker run -v $(pwd)/server:/go/server --name tcp_server -d golang /go/server
func main() {
	if err := epoll.SetLimit(); err != nil {
		panic(err)
	}
	ln, err := net.Listen("tcp", ":8972")
	if err != nil {
		panic(err)
	}
	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()
	epoller, err = epoll.MkEpoll()
	if err != nil {
		panic(err)
	}
	go start()
	for {
		conn, e := ln.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				log.Printf("accept temp err: %v", ne)
				continue
			}
			log.Printf("accept err: %v", e)
			return
		}
		if err := epoller.Add(conn); err != nil {
			log.Printf("failed to add connection %v", err)
			conn.Close()
		}
	}
}

func start() {
	var buf = make([]byte, 8)
	for {
		connections, err := epoller.Wait()
		if err != nil {
			log.Printf("failed to epoll wait %v", err)
			continue
		}
		for _, conn := range connections {
			if conn == nil {
				break
			}
			if _, err := conn.Read(buf); err != nil {
				if err := epoller.Remove(conn); err != nil {
					log.Printf("failed to remove %v", err)
				}
				conn.Close()
			}
		}
	}
}
