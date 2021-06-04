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
