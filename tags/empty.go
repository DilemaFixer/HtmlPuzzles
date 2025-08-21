package tags

import (
	htmlparser "github.com/DilemaFixer/HtmlPuzzles/html"
	"github.com/DilemaFixer/HtmlPuzzles/tools"
)

type EmptyRenderer struct{}

func NewEmptyRenderer() HtmlRenderer {
	return &EmptyRenderer{}
}

func (e *EmptyRenderer) AddChildren(HtmlRenderer) {
	panic("Try add children to EmptyRenderer")
}

func (e *EmptyRenderer) Render(tools.Context) ([]*htmlparser.HtmlTag, error) {
	return make([]*htmlparser.HtmlTag, 0), nil
}

func (e *EmptyRenderer) AddHtml([]*htmlparser.HtmlTag) {
	panic("Try add html to EmptyRenderer")
}
