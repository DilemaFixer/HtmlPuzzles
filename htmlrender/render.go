package htmlrender

type Renderer interface {
	Render(ctx Context) (*RenderedNode, error)
	AddChildren(*Renderer)
}

type RenderedNode struct {
	Html      string
	Childrens []*RenderedNode
}
