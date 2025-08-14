package tagrender

import (
	"fmt"

	render "github.com/DilemaFixer/HtmlPuzzles/htmlrender"
)

type Iterator interface {
	Init(render.Context) error
	IsDone() bool
	GetI() uint64
	Next()
}

type ConstantIterator struct {
	i      uint64
	target uint64
}

func NewConstantIterator(ctx render.Context) (Iterator, error) {
	var value uint64
	var isExist bool

	if value, isExist = ctx.GetUintLayered(render.IterationsCountKey); !isExist {
		return nil, fmt.Errorf("Constant iterator creating error: %s key exist in uint64 map , in layer %d", render.IterationsCountKey, ctx.GetCurrentLayer())
	}

	return &ConstantIterator{
		target: value,
		i:      0,
	}, nil
}

func (c *ConstantIterator) Init(ctx render.Context) error {
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
	childrens []render.Renderer
}

func (f *ForRenderer) Render(ctx render.Context) (*render.RenderedNode, error) {
	currentNode := &render.RenderedNode{
		Childrens: make([]*render.RenderedNode, 0),
	}

	if err := f.iterator.Init(ctx); err != nil {
		return nil, err
	}

	for f.iterator.IsDone() {
		for _, children := range f.childrens {
			node, err := children.Render(ctx)

			if err != nil {
				return nil, err
			}

			currentNode.Childrens = append(currentNode.Childrens, node)
		}
		f.iterator.Next()
	}

	return currentNode, nil
}

func (f *ForRenderer) AddChildren(children render.Renderer) {
	f.childrens = append(f.childrens, children)
}
