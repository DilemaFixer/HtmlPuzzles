package html

import (
	"strconv"
)

type HtmlBuilder struct {
	root    *HtmlTag
	current *HtmlTag
	stack   []*HtmlTag
}

func NewBuilder(tagName string) *HtmlBuilder {
	root := &HtmlTag{
		Name:       tagName,
		Attributes: make(map[string]HtmlAttribute),
	}
	return &HtmlBuilder{
		root:    root,
		current: root,
		stack:   []*HtmlTag{root},
	}
}

func (b *HtmlBuilder) AddTag(name string) *HtmlBuilder {
	tag := &HtmlTag{
		Name:       name,
		Parent:     b.current,
		Attributes: make(map[string]HtmlAttribute),
	}
	b.current.Children = append(b.current.Children, tag)
	b.stack = append(b.stack, tag)
	b.current = tag
	return b
}

func (b *HtmlBuilder) AttrString(name, value string) *HtmlBuilder {
	b.current.SetAttribute(name, value)
	return b
}

func (b *HtmlBuilder) AttrBool(name string, value bool) *HtmlBuilder {
	b.current.SetAttribute(name, strconv.FormatBool(value))
	return b
}

func (b *HtmlBuilder) AttrInt(name string, value int) *HtmlBuilder {
	b.current.SetAttribute(name, strconv.Itoa(value))
	return b
}

func (b *HtmlBuilder) AttrInt64(name string, value int64) *HtmlBuilder {
	b.current.SetAttribute(name, strconv.FormatInt(value, 10))
	return b
}

func (b *HtmlBuilder) AttrUint64(name string, value uint64) *HtmlBuilder {
	b.current.SetAttribute(name, strconv.FormatUint(value, 10))
	return b
}

func (b *HtmlBuilder) AttrFloat32(name string, value float32) *HtmlBuilder {
	b.current.SetAttribute(name, strconv.FormatFloat(float64(value), 'f', -1, 32))
	return b
}

func (b *HtmlBuilder) AttrFloat64(name string, value float64) *HtmlBuilder {
	b.current.SetAttribute(name, strconv.FormatFloat(value, 'f', -1, 64))
	return b
}

func (b *HtmlBuilder) SetText(content string) *HtmlBuilder {
	b.current.InnerContent = content
	return b
}

func (b *HtmlBuilder) SelfClosing() *HtmlBuilder {
	b.current.IsSelfClosing = true
	return b
}

func (b *HtmlBuilder) Up() *HtmlBuilder {
	if len(b.stack) > 1 {
		b.stack = b.stack[:len(b.stack)-1]
		b.current = b.stack[len(b.stack)-1]
	}
	return b
}

func (b *HtmlBuilder) Build() *HtmlTag {
	return b.root
}
