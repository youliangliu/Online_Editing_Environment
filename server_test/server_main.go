package main

import (
	"net"
	"net/http"
	"net/rpc"
)

type HandleSet struct {
}

func (this *HandleSet) Test(in int, res *int) error {
	(*res) = in
	return nil
}

func (this *HandleSet) RequestChange(in int, res *int) error {
	(*res) = in
	return nil
}

func (this *HandleSet) TransferChange(in int, res *int) error {
	(*res) = in
	return nil
}
func ServeBack(ready chan bool) error {
	addr := "localhost:31696"
	server := rpc.NewServer()
	rpc_struct := new(HandleSet)
	server.RegisterName("HandleSet", rpc_struct)
	//rpc.Register(b.Store)
	//rpc.HandleHTTP()
	l, e := net.Listen("tcp", addr)

	ready<-true

	if e != nil {
		return e
	}
	return http.Serve(l, server)
}

func main() {
	ready := make(chan bool)
	go ServeBack(ready)
	<-ready
	for true {
		// Literally  do nothing
	}
}
