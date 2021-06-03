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
func ServeBack() error {
	addr := "localhost:31696"
	server := rpc.NewServer()
	rpc_struct := new(HandleSet)
	server.RegisterName("HandleSet", rpc_struct)
	//rpc.Register(b.Store)
	//rpc.HandleHTTP()
	l, e := net.Listen("tcp", addr)
	if e != nil {
		return e
	}
	return http.Serve(l, server)
}

func main() {
	ServeBack()
}
