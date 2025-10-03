package parsers

import (
	"github.com/DilemaFixer/HtmlPuzzles/examples/nodes"
	"github.com/DilemaFixer/HtmlPuzzles/render"
)

type HostParser struct{}

func NewHostParser() render.NodeParser {
	return &HostParser{}
}

func (h *HostParser) GetTarget() string {
	return "any"
}

func (h *HostParser) Parser(htmlNode *render.HtmlNode, childrenCount uint64) (render.Node, error) {
	return nodes.NewHostNode(htmlNode, childrenCount), nil
}
