package proxyhandler

import (
	"encoding/xml"
	"io"
	"log"
	"strings"
)

// A Parser reads and decodes XML objects from an input stream.
type Parser struct {
	reader io.Reader
	schema map[string]xsdElement
}

type Element struct {
	parent *Element
	node   *Node
	label  string
}

// Creates a New Parser.
func NewParser(r io.Reader, schema map[string]xsdElement) *Parser {
	return &Parser{reader: r, schema: schema}
}

// Decodes XML
// label indicates the element to return, otherwise the root element is returned
func (dec *Parser) Decode(root *Node, label string) (*Node, error) {

	xmlDec := xml.NewDecoder(dec.reader)

	element := &Element{parent: nil, node: root}

	var section *Node = root

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
				node:   &Node{DataType: schema.Type, Complex: schema.isComplex(), List: schema.isList()},
				label:  se.Name.Local,
			}
		case xml.CharData:
			element.node.Data = strings.TrimSpace(string(xml.CharData(se)))
		case xml.EndElement:
			if element.parent != nil {
				element.parent.node.add(element.label, element.node)
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
