package tags

import (
	htmlparser "github.com/DilemaFixer/HtmlPuzzles/html"
	"github.com/DilemaFixer/HtmlPuzzles/tools"
)

type HtmlRenderer interface {
	AddChildren(HtmlRenderer)
	Render(tools.Context) ([]*htmlparser.HtmlTag, error)
	AddHtml([]*htmlparser.HtmlTag)
}
