package render

import (
	"fmt"

	htmlparser "github.com/DilemaFixer/HtmlPuzzles/html"
	"github.com/DilemaFixer/HtmlPuzzles/tags"
)

type RendererParser func(*htmlparser.HtmlTag)

type RenderersParser struct {
	binds   map[string]RendererParser
	builder *RenderersBuilder
}

func NewRenderersParser() *RenderersParser {
	return &RenderersParser{
		binds:   make(map[string]RendererParser),
		builder: NewRenderersBuilderWithDefaultSetUp(),
	}
}

func NewRederersParserWithDefaultSetUp() *RenderersParser {
	p := &RenderersParser{
		binds:   make(map[string]RendererParser),
		builder: NewRenderersBuilderWithDefaultSetUp(),
	}
	p.defaultSetUpForParser()
	return p
}

func (rp *RenderersParser) defaultSetUpForParser() {
	//TODO: when base tags will be exist , add it as default
}

func (rp *RenderersParser) SetCustomBuilder(builder *RenderersBuilder) {
	if builder == nil {
		return
	}

	rp.builder = builder
}

func (rp *RenderersParser) Bind(target string, parser RendererParser) {
	rp.binds[target] = parser
}

func (rp *RenderersParser) HasParserFor(tag string) bool {
	_, exist := rp.binds[tag]
	return exist
}

func (rp *RenderersParser) Parse(tag *htmlparser.HtmlTag) error {
	if !rp.HasParserFor(tag.Name) {
		return fmt.Errorf("Parser for html tag %s not exist in bindings", tag.Name)
	}

	parser := rp.binds[tag.Name]
	parser(tag)
	return nil
}

func (rp *RenderersParser) ParseHtmlToRenderTree(htmlASTs []*htmlparser.HtmlTag) ([]tags.HtmlRenderer, error) {
	if len(htmlASTs) == 0 {
		return nil, fmt.Errorf("Html to render tree parsing error: input html slice is empty or nil")
	}

	result := make([]tags.HtmlRenderer, 0, len(htmlASTs))

	for i, htmlAST := range htmlASTs {
		if htmlAST == nil {
			return nil, fmt.Errorf("Html to render tree parsing error: item in pos %d is nil ptr", i)
		}
		rendererBranch, err := rp.parseHtmlAstToRendererTree(htmlAST)

		if err != nil {
			return nil, fmt.Errorf("Html to render tree parsing error: %s", err.Error())
		}

		result = append(result, rendererBranch)
	}

	return result, nil
}

func (rp *RenderersParser) parseHtmlAstToRendererTree(htmlAst *htmlparser.HtmlTag) (tags.HtmlRenderer, error) {
	if rp.HasParserFor(htmlAst.Name) {
		err := rp.Parse(htmlAst)
		if err != nil {
			return nil, err
		}

		renderer, err := rp.builder.Build(htmlAst)
		if err != nil {
			return nil, err
		}

		renderer.AddHtml(htmlAst.Children)

		if len(htmlAst.Children) > 0 {
			for _, children := range htmlAst.Children {
				childrenRenderer, err := rp.parseHtmlAstToRendererTree(children)
				if err != nil {
					return nil, err
				}

				renderer.AddChildren(childrenRenderer)
			}
		}
		return renderer, nil
	}

	if len(htmlAst.Children) > 0 {
		group := tags.NewGroupRenderer()
		for _, child := range htmlAst.Children {
			childRenderer, err := rp.parseHtmlAstToRendererTree(child)
			if err != nil {
				return nil, err
			}
			if childRenderer != nil {
				group.AddChildren(childRenderer)
			}
		}
		return group, nil
	}

	return tags.NewEmptyRenderer(), nil
}
