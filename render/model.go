package render

import (
	htmlparser "github.com/DilemaFixer/HtmlParser"
	"github.com/DilemaFixer/HtmlPuzzles/tools"
)

type RenderNode interface {
	Render(ctx *tools.Context) ([]*htmlparser.HtmlTag, error)
	AddChildren(node RenderNode)
}
