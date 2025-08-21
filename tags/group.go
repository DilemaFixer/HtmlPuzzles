package tags

import (
	htmlparser "github.com/DilemaFixer/HtmlPuzzles/html"
	"github.com/DilemaFixer/HtmlPuzzles/tools"
)

type GroupRenderer struct {
	group []HtmlRenderer
}

func NewGroupRenderer() *GroupRenderer {
	return &GroupRenderer{
		group: make([]HtmlRenderer, 0),
	}
}

func (g *GroupRenderer) AddChildren(rndrr HtmlRenderer) {
	g.group = append(g.group, rndrr)
}

func (g *GroupRenderer) Render(ctx tools.Context) ([]*htmlparser.HtmlTag, error) {
	htmls := make([]*htmlparser.HtmlTag, 0)

	for _, item := range g.group {
		result, err := item.Render(ctx)
		if err != nil {
			return nil, err
		}
		htmls = append(htmls, result...)
	}
	return htmls, nil
}

func (g *GroupRenderer) AddHtml([]*htmlparser.HtmlTag) {
	panic("Try add html to GroupRenderer")
}
