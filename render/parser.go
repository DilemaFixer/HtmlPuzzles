package render

import (
	"fmt"
	"strings"

	htmlparser "github.com/DilemaFixer/HtmlParser"
)

type TagParser func(*htmlparser.HtmlTag) (RenderNode, error)
type HtmlValidator func(*htmlparser.HtmlTag) error

type TagsRenderer struct {
	binds map[string]TagsParserBind
}

type TagsParserBind struct {
	validator HtmlValidator
	parser    TagParser
}

func NewTagsRenderer() *TagsRenderer {
	return &TagsRenderer{
		binds: make(map[string]TagsParserBind),
	}
}

func (p *TagsRenderer) Bind(key string, parser TagParser, validator HtmlValidator) {
	key = strings.TrimSpace(key)
	if key == "" {
		return
	}

	if parser == nil || validator == nil {
		return
	}

	p.binds[key] = TagsParserBind{
		validator: validator,
		parser:    parser,
	}
}

func (p *TagsRenderer) HasBindFor(key string) bool {
	_, exist := p.binds[key]
	return exist
}

func (p *TagsRenderer) tryRender(htmlNode *htmlparser.HtmlTag) (RenderNode, error) {
	if !p.HasBindFor(htmlNode.Name) {
		return nil, fmt.Errorf("parser for tag '%s' not exist", htmlNode.Name)
	}

	bind := p.binds[htmlNode.Name]
	if err := bind.validator(htmlNode); err != nil {
		return nil, err
	}

	render, err := bind.parser(htmlNode)
	if err != nil {
		return nil, err
	}

	return render, nil
}

func (p *TagsRenderer) HtmlToRenderTree(htmlNodes []*htmlparser.HtmlTag) ([]RenderNode, error) {
	if len(htmlNodes) == 0 {
		return nil, fmt.Errorf("Html to render tree parsing error: input is nil or empyt slice")
	}

	renderNodes := make([]RenderNode, 0, len(htmlNodes))
	for _, htmlNode := range htmlNodes {
		renderBranch, err := p.parseHtml(htmlNode)
		if err != nil {
			return nil, fmt.Errorf("Html to render tree parsing error: %w", err)
		}
		renderNodes = append(renderNodes, renderBranch)
	}

	return renderNodes, nil
}

func (p *TagsRenderer) parseHtml(htmlNode *htmlparser.HtmlTag) (RenderNode, error) {
	if p.HasBindFor(htmlNode.Name) {
		renderer, err := p.tryRender(htmlNode)
		if err != nil {
			return nil, err
		}

		for _, htmlChildren := range htmlNode.Children {
			childrenRenderer, err := p.parseHtml(htmlChildren)
			if err != nil {
				return nil, err
			}
			renderer.AddChildren(childrenRenderer)
		}
		return renderer, nil
	}

	hNode := NewHostNode(htmlNode)

	for _, htmlChildren := range htmlNode.Children {
		childrenRenderer, err := p.parseHtml(htmlChildren)
		if err != nil {
			return nil, err
		}
		hNode.AddChildren(childrenRenderer)
	}
	return hNode, nil
}
