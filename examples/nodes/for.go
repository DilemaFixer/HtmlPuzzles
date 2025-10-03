package nodes

import (
	render "github.com/DilemaFixer/HtmlPuzzles/render"
)

type ForNode struct {
	children  []render.Node
	itr_count uint64
}

func NewForNode(itr_count uint64, childrenCount uint64) render.Node {
	return &ForNode{
		itr_count: itr_count,
		children:  make([]render.Node, 0, childrenCount),
	}
}

func (f *ForNode) Render(ctx *render.Context) (render.RenderResult, error) {
	ctx.LayerUp()
	defer ctx.LayerDown()

	result := render.CompositeResult{}
	for i := uint64(0); i < f.itr_count; i++ {
		for _, subNode := range f.children {
			rendered, err := subNode.Render(ctx)
			if err != nil {
				return nil, err
			}
			result.Children = append(result.Children, rendered)
		}
	}

	return result, nil
}

func (f *ForNode) AddChildren(node render.Node) {
	if node == nil {
		return
	}
	f.children = append(f.children, node)
}
