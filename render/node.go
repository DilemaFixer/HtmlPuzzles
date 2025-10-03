package render

type Node interface {
	Render(ctx *Context) (RenderResult, error)
	AddChildren(node Node)
}
