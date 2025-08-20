package tagrender

import (
	"fmt"

	"github.com/DilemaFixer/HtmlPuzzles/htmlparser"
	"github.com/DilemaFixer/HtmlPuzzles/tools"
)

type Iterator interface {
	Init(tools.Context) error
	IsDone() bool
	GetI() uint64
	Next()
}

type ConstantIterator struct {
	i      uint64
	target uint64
}

func NewConstantIterator(count uint64) (Iterator, error) {
	if count == 0 {
		return nil, fmt.Errorf("Constant iterator creating error: count is zero")
	}

	return &ConstantIterator{
		target: count,
		i:      0,
	}, nil
}

func (c *ConstantIterator) Init(ctx tools.Context) error {
	return nil
}

func (c *ConstantIterator) IsDone() bool {
	return c.i >= c.target
}

func (c *ConstantIterator) GetI() uint64 {
	return c.i
}

func (c *ConstantIterator) Next() {
	c.i++
}

//-----------------------------------------------------------------

type ForRenderer struct {
	iterator  Iterator
	childrens []Renderer
	html      *htmlparser.HtmlTag
}

func NewForRenderer(iterator Iterator) *ForRenderer {
	return &ForRenderer{
		iterator:  iterator,
		childrens: make([]Renderer, 0),
	}
}

func (f *ForRenderer) Render(ctx *tools.Context) (*htmlparser.HtmlTag, error) {
	return nil, nil
}

func (f *ForRenderer) AddChildren(children Renderer) {
	f.childrens = append(f.childrens, children)
}

func (f *ForRenderer) EstablishResponsibilityForRendering(html *htmlparser.HtmlTag) {
	f.html = html
}
