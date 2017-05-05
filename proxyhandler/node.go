package proxyhandler

import (
	"encoding/json"
	"strconv"
)

// Node is a data element on a tree
type Node struct {
	Nodes    map[string]Nodes
	Data     string
	DataType string
	Complex  bool
	List     bool
}

// Nodes is a list of nodes
type Nodes []*Node

// add appends a node to the list of children
func (n *Node) add(label string, node *Node) {
	if n.Nodes == nil {
		n.Nodes = map[string]Nodes{}
	}
	n.Nodes[label] = append(n.Nodes[label], node)
}

// encodes to json
func (n *Node) encode() ([]byte, error) {
	object := n.serialize()
	bytes, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// serializes a node to plain anonymous objects
func (n *Node) serialize() interface{} {
	if n.Complex {
		object := make(map[string]interface{})
		for label, children := range n.Nodes {
			if len(children) > 1 || (len(children) == 1 && children[0].List) {
				s := make([]interface{}, 0)
				for _, c := range children {
					s = append(s, c.serialize())
				}
				object[label] = s
			} else {
				object[label] = children[0].serialize()
			}
		}
		return object
	} else {
		return n.value()
	}
}

// converts to datatype
func (n *Node) value() interface{} {
	var val interface{}
	var err error
	switch n.DataType {
	case "xs:boolean":
		val, err = strconv.ParseBool(n.Data)
	case "xs:integer":
		val, err = strconv.ParseInt(n.Data, 10, 64)
	case "xs:decimal":
		val, err = strconv.ParseFloat(n.Data, 64)
	default:
		val, err = n.Data, nil
	}
	if err != nil {
		return nil
	}
	return val
}
