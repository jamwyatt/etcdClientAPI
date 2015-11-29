// Simple etcd client library to interface with etcd's HTTP API
package etcdMisc

/*
etcdClientAPI is a simple golang library to interface with etcd's API
Copyright (C) 2015 J. Robert Wyatt

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
*/

import (
	"fmt"
)

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
		return fmt.Sprintf("EtcdResponse Action: \"%v\"\nnode:\n%vprevNode:\n%v\n", r.Action, r.Node, r.PrevNode)
	} else {
		return fmt.Sprintf("EtcdResponse ERR: %v", r.err)
	}
}

func (r EtcdResponse) GetError() error {
	return r.err
}
