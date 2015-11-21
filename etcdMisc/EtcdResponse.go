package etcdMisc

import (
	"fmt"
)

type ListNode struct {
	CreatedIndex int
	Key          string
	Value        string
}

type Node struct {
	CreatedIndex  int
	Key           string
	ModifiedIndex int
	Value         string
	Dir           bool
	Nodes         []Node
}

func (n Node) printNode(prefix string) string {
	str := fmt.Sprintf("%s[%v/%v dir=%t - \"%v\"=\"%v\"]", prefix, n.CreatedIndex, n.ModifiedIndex, n.Dir, n.Key, n.Value)
	if n.Dir {
		str += "\n"
		for _, node := range n.Nodes {
			str += node.printNode(prefix + "  ")
			if !node.Dir {
				str += "\n"
			}
		}
	}
	return str
}

func (n Node) String() string {
	return n.printNode("")
}

type EtcdResponse struct {
	Action   string
	Node     Node
	PrevNode Node

	// Next three are only set when the response is an error
	Cause     string
	ErrorCode int
	Message   string

	// hidden from JSON processing, used for error responses
	err error
}

func (r EtcdResponse) String() string {
	if r.err == nil {
		return fmt.Sprintf("EtcdResponse Action: \"%v\" node: %v\nprevNode: %v", r.Action, r.Node, r.PrevNode)
	} else {
		return fmt.Sprintf("EtcdResponse ERR: %v", r.err)
	}
}

func (r EtcdResponse) GetError() error {
	return r.err
}
