package htmlrender

import (
	"github.com/DilemaFixer/HtmlPuzzles/htmlparser"
)

type Renderer interface {
	Render(ctx Context) (*RenderedNode, error)
	AddChildren(*Renderer)
}

type RenderedNode struct {
	Html      string
	Childrens []*RenderedNode
}

func ParseToRenderedNodes(html string) (*RenderedNode, error) {
	htmlTokens, err := htmlparser.ParseHtml(html)

	if err != nil {
		return nil, err
	}

	for _, token := range htmlTokens {
		print(token)
	}

	return nil, nil
}
