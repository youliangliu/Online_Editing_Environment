package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"sync"
	"time"
)
type CrdtArray struct{
	Crdts []Crdt
}
type Clock struct {
	Pos int    `json:"pos"`
	IP  string `json:"IP"`
}

type Crdt struct {
	Value     string `json:"value"`
	Id        []Clock `json:"id"`
	Timestamp string `json:"timestamp"`
	Operation bool `json:"operation"`
}

type Server struct {
	crdts    []Crdt
	backends []string
	index    int
	server_num int
	last_communicate_time []int64
	last_connect_time []int64
	last_try_connect int64
	connected bool
	crdt_lock *sync.Mutex
}

var myserver Server

func (self* Server) sync_crdts(index int) error {
	myserver.last_communicate_time[index] = time.Now().UnixNano()
	if index == myserver.index {
		return nil
	}
	conn, e := rpc.DialHTTP("tcp", self.backends[index])
	if e != nil {
		return e
	}

	// perform the call
	var returned_crdts CrdtArray
	e = conn.Call("Server.CrdtsSync", self.crdts, &returned_crdts)
	if e != nil {
		conn.Close()
		log.Println("Sync Error:", e)
		return e
	}
	conn.Close()
	myserver.last_connect_time[index] = time.Now().UnixNano()
	return nil
}

func (self *Server) Broadcast(command Crdt, ret *bool) error {
	quorum_count:=0
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
			conn.Close()
			myserver.last_connect_time[i] = time.Now().UnixNano()
			quorum_count++
		}
	}
	if quorum_count<myserver.server_num/2 {
		myserver.connected = false
		return errors.New("Disconnected!")
	} else{
		self.PutCommand(command, ret)
		myserver.connected = true
		return nil
	}
}

func (self *Server) PutCommand(command Crdt, ret *bool) error {
	myserver.crdt_lock.Lock()
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
					last := []Crdt{}
					last = append(last, self.crdts[i:]...)
					self.crdts = append(first, command)
					self.crdts = append(self.crdts, last...)
					flag = false
					break
				}
			} else if result > 0 {
				first := self.crdts[:i]
				last := []Crdt{}
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
					last := []Crdt{}
					last = append(last, self.crdts[i:]...)
					self.crdts = append(first, command)
					self.crdts = append(self.crdts, last...)
					flag = false
					break
				}
			} else if result > 0 {
				first := self.crdts[:i]
				last := []Crdt{}
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
	myserver.crdt_lock.Unlock()
	return nil
}

func (self *Server) try_connect() bool {
	myserver.last_try_connect = time.Now().UnixNano()
	quorum_count:=0
	for i := 0; i < len(self.backends); i++ {
		if i != self.index {
			conn, e := rpc.DialHTTP("tcp", self.backends[i])
			if e != nil {
				continue
			}
			quorum_count++
			conn.Close()
		}
	}
	if quorum_count<myserver.server_num/2 {
		myserver.connected = false
		return false
	} else{
		myserver.connected = true
		return true
	}
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

func clock_compare(c1 Clock, c2 Clock) int {
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

func crdt_compare(c1 Crdt, c2 Crdt) int {
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

func (self *Server) crdts_difference(clientCrdts []Crdt) []Crdt {
	var result []Crdt
	for i := 0; i < len(self.crdts); i++ {
		if crdt_valIdation(self.crdts[i], clientCrdts) == true {
			result = append(result, self.crdts[i])
		}
	}
	return result
}

func (self *Server) CrdtsSync(clientCrdts []Crdt, returned_crdts *CrdtArray) error {
	returned_crdts.Crdts = []Crdt{}

	for i := 0; i < len(self.crdts); i++ {
		if crdt_valIdation(self.crdts[i], clientCrdts) == true {
			returned_crdts.Crdts = append(returned_crdts.Crdts, self.crdts[i])
		}
	}

	var append_list []Crdt
	for i := 0; i < len(clientCrdts); i++ {
		if crdt_valIdation(clientCrdts[i], self.crdts) == true {
			append_list = append(append_list, clientCrdts[i])
		}
	}

	var ret bool
	for i := 0; i < len(append_list); i++ {
		self.PutCommand(append_list[i], &ret)
	}

	return nil
}

func crdt_valIdation(command Crdt, crdts []Crdt) bool {
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

func parse_string_to_crdts(b []byte) []Crdt {
	ori := []Crdt{}
	json.Unmarshal(b, &ori)
	return ori
}

func parse_string_to_crdt(b []byte) Crdt {
	var ori Crdt
	json.Unmarshal(b, &ori)
	log.Println(string(b))
	return ori
}

func parse_crdts_to_string(crdts []Crdt) []byte {
	b, _ := json.Marshal(crdts)
	return b
}

func client_put_crdt(w http.ResponseWriter, r *http.Request) {
	if myserver.connected==false {
		w.Write([]byte("Disconnect"))
		return
	}
	// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	switch r.Method {
	case "GET":
		crdt := parse_string_to_crdt([]byte(r.FormValue("content")))
		log.Println("Receive single:", crdt)
		var ret bool
		myserver.Broadcast(crdt, &ret)

	default:
		fmt.Fprintf(w, "Methods error!")
	}
}

func client_update(w http.ResponseWriter, r *http.Request) {
	if myserver.connected==false {
		log.Println("Reject client", myserver.connected)
		w.Write([]byte("Disconnect"))
		return
	}
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	switch r.Method {
	case "GET":
		//if it is a crdt array
		crdts := parse_string_to_crdts([]byte(r.FormValue("content")))
		log.Println("Receive update:", crdts)
		var returned_crdts []Crdt
		returned_crdts = myserver.crdts_difference(crdts)
		var returned_string []byte
		returned_string = parse_crdts_to_string(returned_crdts)
		w.Write(returned_string)

	default:
		fmt.Fprintf(w, "Methods error!")
	}
}


func test_parser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	switch r.Method {
	case "GET":
		crdt := parse_string_to_crdt([]byte(r.FormValue("content")))
		log.Println(crdt)


	default:
		fmt.Fprintf(w, "Methods error!")
	}
}

func test_parser_array(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	switch r.Method {
	case "GET":
		//if it is a crdt array
		crdts := parse_string_to_crdts([]byte(r.FormValue("content")))
		log.Println(crdts)
		returned_string := parse_crdts_to_string(crdts)
		w.Write(returned_string)

	default:
		fmt.Fprintf(w, "Methods error!")
	}
}

func main() {
	var communicate_interval int64
	var try_interval int64
	var connect_interval int64
	port_list := []string{"8080","31267"}
	communicate_interval = 2000000000
	try_interval = 2000000000
	connect_interval = 6000000000
	myserver.server_num = len(port_list)
	myserver.index = 0
	myserver.backends = []string{}
	for i:=0;i!=myserver.server_num;i++ {
		myserver.backends = append(myserver.backends, "localhost:"+port_list[i])
	}
	myserver.connected = false
	myserver.crdt_lock = &sync.Mutex{}
	myserver.last_try_connect = time.Now().UnixNano()
	myserver.last_communicate_time = []int64{}
	for i:=0;i< len(myserver.backends);i++{
		myserver.last_communicate_time = append(myserver.last_communicate_time, int64(0))
		myserver.last_connect_time = append(myserver.last_connect_time, int64(0))
	}



	http.HandleFunc("/send", client_put_crdt)
	http.HandleFunc("/update", client_update)

	//http.HandleFunc("/send", test_parser)
	//http.HandleFunc("/update", test_parser_array)
	rpc.Register(&myserver)
	rpc.HandleHTTP()

	fmt.Printf("Starting server...\n")
	go func() {
		if err := http.ListenAndServe(":"+port_list[myserver.index], nil); err != nil {
			log.Fatal(err)
		}
	}()

	for ;; {
		time.Sleep(time.Duration(300)*time.Millisecond)
		current_time := time.Now().UnixNano()
		if myserver.connected == false && current_time>myserver.last_try_connect+try_interval{
			myserver.try_connect()
		}
		for i:= 0; i<myserver.server_num;i++ {
			if current_time>myserver.last_communicate_time[i]+communicate_interval{
				myserver.sync_crdts(i)
			}
		}
		quorum_count:=0
		for i:= 0; i<myserver.server_num;i++ {
			if current_time<myserver.last_connect_time[i]+connect_interval && i!=myserver.index{
				quorum_count++
			}
		}
		if quorum_count < myserver.server_num/2 {
			myserver.connected = false
		} else {
			myserver.connected = true
		}
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
