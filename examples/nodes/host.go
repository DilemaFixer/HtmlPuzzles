package nodes

import (
	render "github.com/DilemaFixer/HtmlPuzzles/render"
)

type HostNode struct {
	children []render.Node
	hosted   *render.HtmlNode
}

func NewHostNode(hosted *render.HtmlNode, childrenCount uint64) render.Node {
	return &HostNode{
		hosted:   hosted,
		children: make([]render.Node, 0, childrenCount),
	}
}

func (h *HostNode) Render(ctx *render.Context) (render.RenderResult, error) {
	ctx.LayerUp()
	defer ctx.LayerDown()

	composite := render.CompositeResult{}

	for _, child := range h.children {
		res, err := child.Render(ctx)
		if err != nil {
			return nil, err
		}
		composite.Children = append(composite.Children, res)
	}

	return render.HostResult{
		Host:     h.hosted,
		Children: composite,
	}, nil
}

func (h *HostNode) AddChildren(node render.Node) {
	if h.children == nil {
		return
	}
	h.children = append(h.children, node)
}
