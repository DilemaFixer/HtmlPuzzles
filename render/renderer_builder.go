package render

import (
	"fmt"
	"strings"

	htmlparser "github.com/DilemaFixer/HtmlPuzzles/html"
)

type RenderersBuilder struct {
	builders map[string]BuilderBind
}

type BuilderBind struct {
	builder   func(*htmlparser.HtmlTag) HtmlRenderer
	validator func(*htmlparser.HtmlTag) error
}

func NewRenderersBuilder() *RenderersBuilder {
	b := &RenderersBuilder{
		builders: make(map[string]BuilderBind),
	}
	b.defaultSetUpForBuilder()
	return b
}

func (b *RenderersBuilder) defaultSetUpForBuilder() {
	//TODO: when base tags will be exist , add it as default
}

func (b *RenderersBuilder) Bind(target string, builder func(*htmlparser.HtmlTag) HtmlRenderer,
	validator func(*htmlparser.HtmlTag) error) {

	if target == "" || target == " " {
		return
	}

	if builder == nil || validator == nil {
		return
	}

	b.builders[target] = BuilderBind{builder: builder, validator: validator}
}

func (b *RenderersBuilder) Build(tagName string, tag *htmlparser.HtmlTag) (HtmlRenderer, error) {
	tagName = strings.TrimSpace(tagName)
	if tagName == "" || tagName == " " {
		return nil, fmt.Errorf("Renderer building error: target tag name is empyt string")
	}

	bind, exist := b.builders[tagName]

	if !exist {
		return nil, fmt.Errorf("Renderer building error: builder for tag %s not exst", tagName)
	}

	if err := bind.validator(tag); err != nil {
		return nil, fmt.Errorf("Renderer building error: %s", err.Error())
	}

	return bind.builder(tag), nil
}
