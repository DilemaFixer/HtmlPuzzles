package nodes

import "github.com/DilemaFixer/HtmlPuzzles/render"

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

func (f *ForNode) Render(ctx *render.Context) (render.HtmlNodes, error) {
	ctx.LayerUp()
	defer ctx.LayerDown()

	nodes := make(render.HtmlNodes, 0, len(f.children))
	for _, subNode := range f.children {
		rendered, err := subNode.Render(ctx)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, rendered...)
	}

	base := nodes
	for i := uint64(1); i < f.itr_count; i++ {
		nodes = append(nodes, base...)
	}

	return nodes, nil
}

func (f *ForNode) AddChildren(node render.Node) {
	if node == nil {
		return
	}
	f.children = append(f.children, node)
}
