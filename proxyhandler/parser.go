package proxyhandler

import (
	"encoding/xml"
	"io"
	"log"
	"strings"
	"fmt"
)

// A Parser reads XML objects from an input stream.
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

// Parses XML
// label indicates the element to return, otherwise and error is returned
func (p *Parser) Parse(root *Node, label string) (*Node, error) {

	xmlDec := xml.NewDecoder(p.reader)

	element := &Element{parent: nil, node: root}

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
			schema := p.schema[se.Name.Local]
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
	if section == nil {
		return nil, fmt.Errorf("Unable to find section/element with name: %v", label)
	}

	log.Printf("Return element %v\n", root)

	return section, nil
}
