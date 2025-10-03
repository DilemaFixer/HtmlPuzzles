package nodes

import (
	render "github.com/DilemaFixer/HtmlPuzzles/render"
	"github.com/DilemaFixer/HtmlPuzzles/utils"
)

type SyncNode struct {
	children []render.Node
}

func NewSyncNode(childrenCount uint64) render.Node {
	return &SyncNode{
		children: make([]render.Node, 0, childrenCount),
	}
}

func (s *SyncNode) Render(ctx *render.Context) (render.RenderResult, error) {
	fut := utils.NewFuture[render.RenderResult]()

	ctx.Go(func() {
		composite := render.CompositeResult{}
		for _, child := range s.children {
			res, err := child.Render(ctx)
			if err != nil {
				fut.Set(nil, err)
				return
			}
			composite.Children = append(composite.Children, res)
		}

		fut.Set(composite, nil)
	})

	return render.AsyncResult{Future: fut}, nil
}

func (s *SyncNode) AddChildren(node render.Node) {
	if s.children == nil {
		return
	}
	s.children = append(s.children, node)
}
