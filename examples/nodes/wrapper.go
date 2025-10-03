package nodes

import (
	"fmt"

	"github.com/DilemaFixer/HtmlPuzzles/html"
	render "github.com/DilemaFixer/HtmlPuzzles/render"
)

const (
	ErrWrapperSelfClosed = "'wrapper' tag is self-closed, children must not exist"
	ErrUnknownRenderType = "unknown render type"
)

type WrapperNode struct {
	wrapped      string
	attrs        map[string]render.Value
	isSelfClosed bool
	children     []render.Node
}

func NewWrapperNode(childrenCount uint64, wrapped string, isSelfClosed bool, attrs map[string]render.Value) render.Node {
	return &WrapperNode{
		wrapped:      wrapped,
		attrs:        attrs,
		isSelfClosed: isSelfClosed,
		children:     make([]render.Node, 0, childrenCount),
	}
}

func (w *WrapperNode) Render(ctx *render.Context) (render.RenderResult, error) {
	htmlBuilder := html.NewBuilder(w.wrapped)
	if w.isSelfClosed {
		htmlBuilder.SelfClosing()
	}

	for key, val := range w.attrs {
		if err := w.applyAttr(htmlBuilder, key, val, ctx); err != nil {
			return nil, err
		}
	}

	composite := render.CompositeResult{}
	for _, child := range w.children {
		result, err := child.Render(ctx)
		if err != nil {
			return nil, err
		}
		composite.Children = append(composite.Children, result)
	}

	return &render.HostResult{
		Host:     htmlBuilder.Build(),
		Children: composite,
	}, nil
}

func (w *WrapperNode) applyAttr(builder *html.HtmlBuilder, key string, val render.Value, ctx *render.Context) error {
	value, valueType, err := val.GetValue(ctx)
	if err != nil {
		return err
	}

	switch valueType {
	case render.Integer:
		builder.AttrInt64(key, value.(int64))
	case render.Uint:
		builder.AttrUint64(key, value.(uint64))
	case render.Float:
		builder.AttrFloat64(key, value.(float64))
	case render.String:
		builder.AttrString(key, value.(string))
	case render.Boolean:
		builder.AttrBool(key, value.(bool))
	default:
		return fmt.Errorf("%s: %v", ErrUnknownRenderType, valueType)
	}
	return nil
}

func (w *WrapperNode) AddChildren(node render.Node) {
	if w.isSelfClosed {
		panic(ErrWrapperSelfClosed)
	}
	if node == nil {
		return
	}
	w.children = append(w.children, node)
}
