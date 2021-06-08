package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
)

type clock struct {
	Pos int    `json:"Pos"`
	IP  string `json:"IP"`
}

type crdt struct {
	Value     byte
	Id        []clock
	Timestamp string
	Operation bool
}

type Server struct {
	crdts    []crdt
	backends []string
	index    int
}

var myserver Server

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
	if command.Operation == true {
		flag := true
		for i := 0; i < len(self.crdts); i++ {
			result := crdt_compare(self.crdts[i], command)
			if result == 0 {
				if command.Timestamp == self.crdts[i].Timestamp {
					flag = false
					break
				}
				if command.Timestamp < self.crdts[i].Timestamp {
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
				if command.Timestamp == self.crdts[i].Timestamp {
					self.crdts[i] = command
					flag = false
					break
				}
				if command.Timestamp < self.crdts[i].Timestamp {
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
		if self.crdts[i].Operation {
			result = result + string(self.crdts[i].Value)
		}
	}
	fmt.Println(result)
	return nil
}

func clock_compare(c1 clock, c2 clock) int {
	if c1.Pos < c2.Pos {
		return -1
	} else if c1.Pos > c2.Pos {
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
	if len(c1.Id) <= len(c2.Id) {
		compareLength = len(c1.Id)
	} else {
		compareLength = len(c2.Id)
	}

	for i := 0; i < compareLength; i++ {
		result := clock_compare(c1.Id[i], c2.Id[i])
		if result != 0 {
			return result
		}
	}

	if len(c1.Id) < len(c2.Id) {
		return -1
	} else if len(c1.Id) > len(c2.Id) {
		return 1
	} else {
		return 0
	}
}

func (self *Server) crdts_difference(clientCrdts []crdt) []crdt {
	var result []crdt
	for i := 0; i < len(self.crdts); i++ {
		if crdt_valIdation(self.crdts[i], clientCrdts) == true {
			result = append(result, self.crdts[i])
		}
	}
	return result
}

func crdt_valIdation(command crdt, crdts []crdt) bool {
	for i := 0; i < len(crdts); i++ {
		result := crdt_compare(crdts[i], command)
		if result == 0 {
			if command.Timestamp == crdts[i].Timestamp {
				if command.Operation == false && crdts[i].Operation == true {
					return true
				} else {
					return false
				}
			}
		}
	}
	return true
}

func parse_string_to_crdts(b []byte) []crdt {
	ori := []crdt{}
	json.Unmarshal(b, &ori)
	return ori
}

func parse_string_to_crdt(b []byte) crdt {
	var ori crdt
	json.Unmarshal(b, &ori)
	return ori
}

func parse_crdts_to_string(crdts []crdt) []byte {
	b, _ := json.Marshal(crdts)
	return b
}

func http_serve_crdt(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	switch r.Method {
	case "GET":
		a := "name:"
		w.Write([]byte(a))
		crdt := parse_string_to_crdt([]byte(r.FormValue("name")))
		var ret bool
		myserver.Broadcast(crdt, &ret)

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func http_serve_crdts(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	switch r.Method {
	case "GET":
		a := "name:"
		w.Write([]byte(a))

		//if it is a crdt array
		crdts := parse_string_to_crdts([]byte(r.FormValue("name")))
		var returned_crdts []crdt
		returned_crdts = myserver.crdts_difference(crdts)
		var returned_string []byte
		returned_string = parse_crdts_to_string(returned_crdts)
		w.Write(returned_string)

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	http.HandleFunc("/http_serve_crdt", http_serve_crdt)
	http.HandleFunc("/http_serve_crdts", http_serve_crdts)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

/*
func main() {
	var a1 clock
	var a2 clock
	var a3 clock
	var a4 clock
	var a5 clock

	a1.Pos = 1
	a1.IP = "A"
	a2.Pos = 2
	a2.IP = "A"
	a3.Pos = 3
	a3.IP = "A"
	a4.Pos = 4
	a4.IP = "A"
	a5.Pos = 5
	a5.IP = "A"

	var b1 clock
	var b2 clock

	b1.Pos = 1
	b1.IP = "B"
	b2.Pos = 2
	b2.IP = "B"

	var c1, c2, c3, c4, c5, c6, c7, c8, c9 crdt
	c1.Id = []clock{}
	c2.Id = []clock{}
	c3.Id = []clock{}
	c4.Id = []clock{}
	c5.Id = []clock{}
	c6.Id = []clock{}
	c7.Id = []clock{}
	c8.Id = []clock{}
	c9.Id = []clock{}

	c1.Id = append(c1.Id, a1)
	c2.Id = append(c2.Id, a2)
	c3.Id = append(c3.Id, a3)
	c4.Id = append(c4.Id, a4)
	c5.Id = append(c5.Id, a4)
	c5.Id = append(c5.Id, b1)
	c6.Id = append(c6.Id, a4)
	c6.Id = append(c6.Id, b2)
	c7.Id = append(c7.Id, a5)

	c8.Id = append(c8.Id, a4)
	c9.Id = append(c9.Id, a4)
	c9.Id = append(c9.Id, b1)

	c1.Operation = true
	c2.Operation = true
	c3.Operation = true
	c4.Operation = true
	c5.Operation = true
	c6.Operation = true
	c7.Operation = true

	c8.Operation = false
	c9.Operation = false

	c1.Value = 'a'
	c2.Value = 'b'
	c3.Value = 'c'
	c4.Value = 'd'
	c5.Value = 'e'
	c6.Value = 'f'
	c7.Value = 'g'

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

	b := parse_crdts_to_string(atext.crdts)
	m := parse_string_to_crdts(b)
	fmt.Println(m[0])
}
*/
