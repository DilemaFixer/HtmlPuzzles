package render

import (
	"fmt"

	htmlparser "github.com/DilemaFixer/HtmlPuzzles/html"
	"github.com/DilemaFixer/HtmlPuzzles/tags"
)

type RenderersBuilder struct {
	builders map[string]BuilderBind
}

type BuilderBind struct {
	builder   func(*htmlparser.HtmlTag) tags.HtmlRenderer
	validator func(*htmlparser.HtmlTag) error
}

func NewRenderersBuilderWithDefaultSetUp() *RenderersBuilder {
	b := &RenderersBuilder{
		builders: make(map[string]BuilderBind),
	}
	b.defaultSetUpForBuilder()
	return b
}

func NewRenderersBuilder() *RenderersBuilder {
	b := &RenderersBuilder{
		builders: make(map[string]BuilderBind),
	}
	return b
}

func (b *RenderersBuilder) defaultSetUpForBuilder() {
	//TODO: when base tags will be exist , add it as default
}

func (b *RenderersBuilder) Bind(target string, builder func(*htmlparser.HtmlTag) tags.HtmlRenderer,
	validator func(*htmlparser.HtmlTag) error) {

	if target == "" || target == " " {
		return
	}

	if builder == nil || validator == nil {
		return
	}

	b.builders[target] = BuilderBind{builder: builder, validator: validator}
}

func (b *RenderersBuilder) Build(tag *htmlparser.HtmlTag) (tags.HtmlRenderer, error) {
	bind, exist := b.builders[tag.Name]

	if !exist {
		return nil, fmt.Errorf("Renderer building error: builder for tag %s not exst", tag.Name)
	}

	if err := bind.validator(tag); err != nil {
		return nil, fmt.Errorf("Renderer building error: %s", err.Error())
	}

	return bind.builder(tag), nil
}

func (b *RenderersBuilder) HasBuilderFor(tagName string) bool {
	_, exist := b.builders[tagName]
	return exist
}
