package parsers

import (
	"fmt"

	"github.com/DilemaFixer/HtmlPuzzles/nodes"
	"github.com/DilemaFixer/HtmlPuzzles/render"
)

type ForParser struct{}

func NewForParser() render.NodeParser {
	return &ForParser{}
}

func (f *ForParser) GetTarget() string {
	return "for"
}

func (f *ForParser) Parser(htmlNode *render.HtmlNode, childrenCount uint64) (render.Node, error) {
	if err := validateNode(htmlNode); err != nil {
		return nil, err
	}

	attr := htmlNode.GetAttribute("itr_count")
	itr_count, err := attr.AsUint64()
	if err != nil {
		return nil, err
	}
	return nodes.NewForNode(itr_count, childrenCount), nil
}

func validateNode(htmlNode *render.HtmlNode) error {
	if !htmlNode.HasAttribute("itr_count") {
		return fmt.Errorf("parsing 'for' html tag error: expected attribute 'itr_count' but it not exists")
	}
	return nil
}
