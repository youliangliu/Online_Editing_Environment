package client

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

type Client struct {
	crdts    []crdt
	backends []string
}

// putcommand

func generate_crdt(input string, crdts []crdt) crdt {
	for i := 0; i < len(crdts); i++ {
		if crdts[i].operation == false {
			crdts = append(crdts[:i], crdts[i+1:]...)
			i--
		}
	}

}
