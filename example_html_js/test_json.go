package main

import (
	"encoding/json"
	"log"
)


type clock struct {
	Pos int `json:"Pos"`
	IP  string `json:"IP"`
}

type crdt struct {
	Value     byte
	Id        []clock
	Timestamp string
	Operation bool
}

func main() {
	t1 := clock{}
	t1.IP = "ip1"
	t1.Pos = 1

	t2 := clock{}
	t2.IP = "ip2"
	t2.Pos = 2

	c:=[]clock{}
	c = append(c, t1)
	c = append(c, t2)

	crdts := []crdt{}

	crdt1 := crdt{
		Value:     1,
		Id:        c,
		Timestamp: "1",
		Operation: false,
	}
	crdt2 := crdt{
		Value:     2,
		Id:        c,
		Timestamp: "2",
		Operation: false,
	}
	crdts = append(crdts, crdt1)
	crdts = append(crdts, crdt2)

	b,_ := json.Marshal(&crdts)
	log.Println(string(b))
	ori:=[]crdt{}
	json.Unmarshal(b,&ori)
	log.Println(ori)
}