package main

import (
	"fmt"
	"log"
	"net/rpc"
)

func main() {
	conn, err := rpc.DialHTTP("tcp", "localhost:31696")
	if err != nil {
		log.Fatalln("dailing error: ", err)
	}

	res := 14
	err = conn.Call("HandleSet.Test", 12, &res)
	if err != nil {
		log.Fatalln("error: ", err)
	}
	fmt.Printf("%d\n", res)

}
