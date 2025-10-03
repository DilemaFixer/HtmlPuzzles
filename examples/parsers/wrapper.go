package parsers

import (
	"fmt"
	"strings"

	"github.com/DilemaFixer/HtmlPuzzles/examples/nodes"
	"github.com/DilemaFixer/HtmlPuzzles/render"
)

const (
	TargetWrapperTag       = "wrapper"
	AttrWrapped            = "wrapped"
	AttrDynamicPrefix      = ":"
	AttrValueSeparator     = "|"
	ErrFormatExpected      = "'path.to.field|type'"
	ErrWrappedAttrNotFound = "expected 'wrapped' attribute not found"
)

type WrapperParser struct{}

func NewWrapperParser() render.NodeParser {
	return &WrapperParser{}
}

func (w *WrapperParser) GetTarget() string {
	return TargetWrapperTag
}

func (w *WrapperParser) Parser(htmlNode *render.HtmlNode, childrenCount uint64) (render.Node, error) {
	if err := validateHtmlNode(htmlNode); err != nil {
		return nil, fmt.Errorf("wrapper parsing failed: %w", err)
	}

	wrappedAttr := htmlNode.GetAttribute(AttrWrapped)
	htmlNode.RemoveAttribute(AttrWrapped)

	values, err := parseAttributes(htmlNode)
	if err != nil {
		return nil, fmt.Errorf("wrapper attribute parsing failed: %w", err)
	}

	return nodes.NewWrapperNode(childrenCount, wrappedAttr.Value, htmlNode.IsSelfClosing, values), nil
}

func parseAttributes(htmlNode *render.HtmlNode) (map[string]render.Value, error) {
	attrs := htmlNode.Attributes
	values := make(map[string]render.Value, len(attrs))

	for key, attr := range attrs {
		if strings.HasPrefix(key, AttrDynamicPrefix) {
			keyWithoutPrefix := strings.TrimPrefix(key, AttrDynamicPrefix)

			if !attr.IsValueExist {
				return nil, fmt.Errorf("dynamic attribute %q has no value", key)
			}

			if err := validateAttributeValue(attr.Value); err != nil {
				return nil, fmt.Errorf("dynamic attribute %q invalid: %w, expected format %s", key, err, ErrFormatExpected)
			}

			typeAndWay := strings.Split(attr.Value, AttrValueSeparator)
			way := strings.Split(typeAndWay[0], ".")
			typ := typeAndWay[1]

			values[keyWithoutPrefix] = render.NewValueFromContext(way, typ)
		} else {
			values[key] = render.NewStaticValue(attr.Value)
		}
	}

	return values, nil
}

func validateAttributeValue(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("value is empty")
	}
	if !strings.Contains(value, AttrValueSeparator) {
		return fmt.Errorf("value must contain %q", AttrValueSeparator)
	}

	typeAndWay := strings.Split(value, AttrValueSeparator)
	if len(typeAndWay) != 2 {
		return fmt.Errorf("invalid format, must be %s", ErrFormatExpected)
	}

	if strings.TrimSpace(typeAndWay[0]) == "" {
		return fmt.Errorf("path is empty")
	}
	if strings.TrimSpace(typeAndWay[1]) == "" {
		return fmt.Errorf("type is empty")
	}

	return nil
}

func validateHtmlNode(htmlNode *render.HtmlNode) error {
	if !htmlNode.HasAttribute(AttrWrapped) {
		return fmt.Errorf(ErrWrappedAttrNotFound)
	}

	attr := htmlNode.GetAttribute(AttrWrapped)
	if strings.TrimSpace(attr.Value) == "" {
		return fmt.Errorf("'wrapped' attribute is empty, expected tag name like <h1>, <div>, <i> etc.")
	}

	return nil
}
