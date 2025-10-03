package parsers

import (
	"github.com/DilemaFixer/HtmlPuzzles/examples/nodes"
	"github.com/DilemaFixer/HtmlPuzzles/render"
)

type SyncParser struct{}

func NewSyncParser() render.NodeParser {
	return &SyncParser{}
}

func (p *SyncParser) GetTarget() string {
	return "sync"
}

func (p *SyncParser) Parser(htmlNode *render.HtmlNode, childrenCount uint64) (render.Node, error) {
	return nodes.NewSyncNode(childrenCount), nil
}
