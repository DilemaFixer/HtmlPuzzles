package nodes

import "github.com/DilemaFixer/HtmlPuzzles/render"

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

func (h *HostNode) Render(ctx *render.Context) (render.HtmlNodes, error) {
	ctx.LayerUp()
	hCopy, err := h.hosted.CloneDown(1)
	if err != nil {
		return nil, err
	}

	hCopy.Children = make(render.HtmlNodes, 0, len(h.children))
	for _, child := range h.children {
		subNodes, err := child.Render(ctx)
		if err != nil {
			return nil, err
		}
		hCopy.Children = append(hCopy.Children, subNodes...)
	}
	ctx.LayerDown()
	return render.HtmlNodes{hCopy}, nil
}

func (h *HostNode) AddChildren(node render.Node) {
	if h.children == nil {
		return
	}
	h.children = append(h.children, node)
}
