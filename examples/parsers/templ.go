package parsers

import (
	"fmt"

	"github.com/DilemaFixer/HtmlPuzzles/examples/nodes"
	"github.com/DilemaFixer/HtmlPuzzles/render"
)

const (
	TargetTemplateTag = "templ"
	AttrSource        = "source"

	ErrTemplateMustBeSelfClosed = "templ tag must be self-closed"
	ErrSourceAttrNotFound       = "templ tag requires 'source' attribute"
)

type TemplateParser struct{}

func NewTemplateParser() render.NodeParser {
	return &TemplateParser{}
}

func (t TemplateParser) GetTarget() string {
	return TargetTemplateTag
}

func (t TemplateParser) Parser(htmlNode *render.HtmlNode, childrenCount uint64) (render.Node, error) {
	if !htmlNode.IsSelfClosing {
		return nil, fmt.Errorf("parsing <%s> error: %s", TargetTemplateTag, ErrTemplateMustBeSelfClosed)
	}

	if !htmlNode.HasAttribute(AttrSource) {
		return nil, fmt.Errorf("parsing <%s> error: %s", TargetTemplateTag, ErrSourceAttrNotFound)
	}

	attr := htmlNode.GetAttribute(AttrSource)
	templPath := attr.AsString()

	return nodes.NewTemplateNode(templPath), nil
}
