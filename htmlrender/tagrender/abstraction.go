package tagrender

import (
	"github.com/DilemaFixer/HtmlPuzzles/htmlparser"
	"github.com/DilemaFixer/HtmlPuzzles/tools"
)

type Renderer interface {
	Render(ctx *tools.Context) (*htmlparser.HtmlTag, error)
	AddChildren(Renderer)
	EstablishResponsibilityForRendering(*htmlparser.HtmlTag)
}
