package etcdMisc

import (
	"fmt"
)

type Node struct {
	CreatedIndex  int
	Key           string
	ModifiedIndex int
	Value         string
}

func (n Node) String() string {
	return fmt.Sprintf("[%v - %v - %v - \"%v\"]", n.CreatedIndex, n.Key, n.ModifiedIndex, n.Value)
}

type EtcdResponse struct {
	Action   string
	Node     Node
	PrevNode Node
	err      error // hidden from JSON processing, used for error responses
}

func (r EtcdResponse) String() string {
	if r.err == nil {
		return fmt.Sprintf("EtcdResponse Action: %v node: %v prevNode: %v", r.Action, r.Node, r.PrevNode)
	} else {
		return fmt.Sprintf("EtcdResponse ERR: %v", r.err)
	}
}

func (r EtcdResponse) GetError() error {
	return r.err
}
