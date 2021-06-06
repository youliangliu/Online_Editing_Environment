package main

import (
	"encoding/json"
	"fmt"
	"net/rpc"
)

type clock struct {
	pos int
	IP  string
}

type crdt struct {
	value     byte
	id        []clock
	timestamp string
	operation bool
}

type Server struct {
	crdts    []crdt
	backends []string
	index    int
}

func (self *Server) Broadcast(command crdt, ret *bool) error {
	self.PutCommand(command, ret)
	for i := 0; i < len(self.backends); i++ {
		if i != self.index {
			conn, e := rpc.DialHTTP("tcp", self.backends[i])
			if e != nil {
				continue
			}

			// perform the call
			e = conn.Call("Server.PutCommand", command, ret)
			if e != nil {
				conn.Close()
				continue
			}
		}
	}
	return nil
}

func (self *Server) PutCommand(command crdt, ret *bool) error {
	if command.operation == true {
		flag := true
		for i := 0; i < len(self.crdts); i++ {
			result := crdt_compare(self.crdts[i], command)
			if result == 0 {
				if command.timestamp == self.crdts[i].timestamp {
					flag = false
					break
				}
				if command.timestamp < self.crdts[i].timestamp {
					first := self.crdts[:i]
					last := []crdt{}
					last = append(last, self.crdts[i:]...)
					self.crdts = append(first, command)
					self.crdts = append(self.crdts, last...)
					flag = false
					break
				}
			} else if result > 0 {
				first := self.crdts[:i]
				last := []crdt{}
				last = append(last, self.crdts[i:]...)
				self.crdts = append(first, command)
				self.crdts = append(self.crdts, last...)
				flag = false
				break
			}
		}
		if flag {
			self.crdts = append(self.crdts, command)
			//fmt.Println(self.crdts)
		}
	} else {
		flag := true
		for i := 0; i < len(self.crdts); i++ {
			result := crdt_compare(self.crdts[i], command)
			if result == 0 {
				if command.timestamp == self.crdts[i].timestamp {
					self.crdts[i] = command
					flag = false
					break
				}
				if command.timestamp < self.crdts[i].timestamp {
					first := self.crdts[:i]
					last := []crdt{}
					last = append(last, self.crdts[i:]...)
					self.crdts = append(first, command)
					self.crdts = append(self.crdts, last...)
					flag = false
					break
				}
			} else if result > 0 {
				first := self.crdts[:i]
				last := []crdt{}
				last = append(last, self.crdts[i:]...)
				self.crdts = append(first, command)
				self.crdts = append(self.crdts, last...)
				flag = false
				break
			}
		}
		if flag {
			self.crdts = append(self.crdts, command)
		}
	}

	return nil
}

func (self *Server) show() error {
	result := ""
	for i := 0; i < len(self.crdts); i++ {
		if self.crdts[i].operation {
			result = result + string(self.crdts[i].value)
		}
	}
	fmt.Println(result)
	return nil
}

func clock_compare(c1 clock, c2 clock) int {
	if c1.pos < c2.pos {
		return -1
	} else if c1.pos > c2.pos {
		return 1
	} else {
		if c1.IP < c2.IP {
			return -1
		} else if c1.IP > c2.IP {
			return 1
		} else {
			return 0
		}
	}
}

func crdt_compare(c1 crdt, c2 crdt) int {
	var compareLength int
	//compareLength = len(c1.clock) < len(c2.clock) ? len(c1.clock) : len(c2.clock)
	if len(c1.id) <= len(c2.id) {
		compareLength = len(c1.id)
	} else {
		compareLength = len(c2.id)
	}

	for i := 0; i < compareLength; i++ {
		result := clock_compare(c1.id[i], c2.id[i])
		if result != 0 {
			return result
		}
	}

	if len(c1.id) < len(c2.id) {
		return -1
	} else if len(c1.id) > len(c2.id) {
		return 1
	} else {
		return 0
	}
}

func (self *Server) crdts_difference(clientCrdts []crdt) []crdt {
	var result []crdt
	for i := 0; i < len(self.crdts); i++ {
		if crdt_validation(self.crdts[i], clientCrdts) == true {
			result = append(result, self.crdts[i])
		}
	}
	return result
}

func crdt_validation(command crdt, crdts []crdt) bool {
	for i := 0; i < len(crdts); i++ {
		result := crdt_compare(crdts[i], command)
		if result == 0 {
			if command.timestamp == crdts[i].timestamp {
				if command.operation == false && crdts[i].operation == true {
					return true
				} else {
					return false
				}
			}
		}
	}
	return true
}

func parseToCrdts(b []byte) []crdt {

}

func main() {
	var a1 clock
	var a2 clock
	var a3 clock
	var a4 clock
	var a5 clock

	a1.pos = 1
	a1.IP = "A"
	a2.pos = 2
	a2.IP = "A"
	a3.pos = 3
	a3.IP = "A"
	a4.pos = 4
	a4.IP = "A"
	a5.pos = 5
	a5.IP = "A"

	var b1 clock
	var b2 clock

	b1.pos = 1
	b1.IP = "B"
	b2.pos = 2
	b2.IP = "B"

	var c1, c2, c3, c4, c5, c6, c7, c8, c9 crdt
	c1.id = []clock{}
	c2.id = []clock{}
	c3.id = []clock{}
	c4.id = []clock{}
	c5.id = []clock{}
	c6.id = []clock{}
	c7.id = []clock{}
	c8.id = []clock{}
	c9.id = []clock{}

	c1.id = append(c1.id, a1)
	c2.id = append(c2.id, a2)
	c3.id = append(c3.id, a3)
	c4.id = append(c4.id, a4)
	c5.id = append(c5.id, a4)
	c5.id = append(c5.id, b1)
	c6.id = append(c6.id, a4)
	c6.id = append(c6.id, b2)
	c7.id = append(c7.id, a5)

	c8.id = append(c8.id, a4)
	c9.id = append(c9.id, a4)
	c9.id = append(c9.id, b1)

	c1.operation = true
	c2.operation = true
	c3.operation = true
	c4.operation = true
	c5.operation = true
	c6.operation = true
	c7.operation = true

	c8.operation = false
	c9.operation = false

	c1.value = 'a'
	c2.value = 'b'
	c3.value = 'c'
	c4.value = 'd'
	c5.value = 'e'
	c6.value = 'f'
	c7.value = 'g'

	var atext Server
	ret := true
	atext.PutCommand(c8, &ret)
	atext.PutCommand(c1, &ret)
	atext.PutCommand(c3, &ret)
	atext.PutCommand(c5, &ret)
	atext.PutCommand(c7, &ret)
	atext.PutCommand(c9, &ret)
	atext.PutCommand(c2, &ret)
	atext.PutCommand(c4, &ret)
	atext.PutCommand(c6, &ret)

	atext.show()

	//a, _ := json.Marshal(atext)
	fmt.Println(atext)

	b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
	var f interface{}
	json.Unmarshal(b, &f)
	m := f.(map[string]interface{})

	test := m["Age"]

}
