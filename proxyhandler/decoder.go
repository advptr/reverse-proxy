package proxyhandler

import (
	"io"
	"encoding/xml"
	"strings"
	"encoding/json"
	"log"
)

// A Decoder reads and decodes XML objects from an input stream.
type Decoder struct {
	reader io.Reader
	schema map[string]xsdElement
}

// Node is a data element on a tree
type Node struct {
	Children map[string]Nodes
	Data     string
	Complex  bool
	List     bool
}

// Nodes is a list of nodes
type Nodes []*Node

// addChild appends a node to the list of children
func (n *Node) addChild(s string, c *Node) {
	if n.Children == nil {
		n.Children = map[string]Nodes{}
	}
	n.Children[s] = append(n.Children[s], c)
}

type Element struct {
	parent   *Element
	node     *Node
	label    string
}


// Creates a New Decoder.
func NewDecoder(r io.Reader, schema map[string]xsdElement) *Decoder {
	return &Decoder{reader: r, schema: schema}
}

// Decodes XML
// returnElementName indicates the element to return, otherwise the root element is returned
func (dec *Decoder) Decode(root *Node, label string) (*Node, error) {

	xmlDec := xml.NewDecoder(dec.reader)

	element := &Element{parent:nil, node:root}

	var section *Node = nil

	for {
		t, err := xmlDec.Token()
		if t == nil || err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		switch se := t.(type) {
		case xml.StartElement:
			schema := dec.schema[se.Name.Local]
			element = &Element{
				parent: element,
				node:   &Node{Complex:schema.isComplex(), List:schema.isList()},
				label:  se.Name.Local,
			}
		case xml.CharData:
			element.node.Data = strings.TrimSpace(string(xml.CharData(se)))
		case xml.EndElement:
			if element.parent != nil {
				element.parent.node.addChild(element.label, element.node)
			}
			if element.label == label {
				section = element.node
			}
			element = element.parent
		}
	}

	log.Printf("Return element %v\n", root)

	return section, nil
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
func (n *Node) serialize() (interface{}) {
	if n.Complex {
		object := make(map[string]interface{})
		for label, children := range n.Children {
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
		return n.Data
	}
}
