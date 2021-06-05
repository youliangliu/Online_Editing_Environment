package main

import "fmt"

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

type text struct {
	crdts []crdt
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

func (self *text) PutCommand(command crdt) error {
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

func (self *text) show() error {
	result := ""
	for i := 0; i < len(self.crdts); i++ {
		if self.crdts[i].operation {
			result = result + string(self.crdts[i].value)
		}
	}
	fmt.Println(result)
	return nil
}

func generate_crdt(crdts []crdt) []crdt {
	for i := 0; i < len(crdts); i++ {
		if crdts[i].operation == false {
			crdts = append(crdts[:i], crdts[i+1:]...)
			i--
		}
	}
	return crdts
}

func show_all_crdt(crdts []crdt) {
	for i := 0; i < len(crdts); i++ {
		fmt.Println(crdts[i])
	}
}

func test(a clock) int {
	a.pos = a.pos + 1
	return a.pos
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

	var atext text
	atext.PutCommand(c8)
	atext.PutCommand(c1)
	atext.PutCommand(c3)
	atext.PutCommand(c5)
	atext.PutCommand(c7)
	atext.PutCommand(c9)
	atext.PutCommand(c2)
	atext.PutCommand(c4)
	atext.PutCommand(c6)

	// atext.show()
	var test_text text
	test_text.crdts = generate_crdt(atext.crdts)
	show_all_crdt(test_text.crdts)

	var testClock clock
	testClock.pos = 1

	//fmt.Println(crdt_compare(c6, c4))

	// var c11 crdt
	// var c22 crdt

	// c11.id = []clock{}
	// c22.id = []clock{}

	// c11.id = append(c11.id, c2)
	// c11.id = append(c11.id, c1)
	// c22.id = append(c22.id, c2)
	// c22.id = append(c22.id, c1)

	// fmt.Println(crdt_compare(c11, c22))

	// array1 := []int{1, 2, 3, 4}
	// array2 := []int{6, 7, 8, 9}
	// array1 = append(array1, array2...)

	// fmt.Println(array1)

}
