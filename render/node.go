package render

type Node interface {
	Render(ctx *Context) (HtmlNodes, error)
	AddChildren(node Node)
}
