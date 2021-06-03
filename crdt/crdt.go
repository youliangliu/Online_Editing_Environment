package crdt

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
					last := self.crdts[i:]
					self.crdts = append(first, command)
					self.crdts = append(self.crdts, last...)
					flag = false
					break
				}
			} else if result > 0 {
				first := self.crdts[:i]
				last := self.crdts[i:]
				self.crdts = append(first, command)
				self.crdts = append(self.crdts, last...)
				flag = false
				break
			}
		}
		if flag {
			self.crdts = append(self.crdts, command)
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
					last := self.crdts[i:]
					self.crdts = append(first, command)
					self.crdts = append(self.crdts, last...)
					flag = false

					break
				}
			} else if result > 0 {
				first := self.crdts[:i]
				last := self.crdts[i:]
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
