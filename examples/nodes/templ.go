package nodes

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/DilemaFixer/HtmlPuzzles/html"
	render "github.com/DilemaFixer/HtmlPuzzles/render"
)

type TemplateNode struct {
	path string
}

func NewTemplateNode(path string) render.Node {
	return &TemplateNode{
		path: path,
	}
}

func (t *TemplateNode) Render(ctx *render.Context) (render.RenderResult, error) {
	ctx.LayerUp()
	defer ctx.LayerDown()

	fullPath := filepath.Join(ctx.GetTemplateRoot(), t.path)

	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, fmt.Errorf("template render error: file %s not found", fullPath)
	}

	if filepath.Ext(fullPath) != ".html" {
		return nil, fmt.Errorf("template render error: file %s must have .html extension", fullPath)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("template render error: %s is a directory, not a file", fullPath)
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("template render error: failed to read %s: %v", fullPath, err)
	}
	source := string(data)
	htmlAsts, err := html.ParseHtml(source)
	if err != nil {
		return nil, fmt.Errorf("source render to html error: failed to parse %s: %v", fullPath, err)
	}

	renderTree, err := ctx.Builder.Build(htmlAsts)
	if err != nil {
		return nil, fmt.Errorf("build render tree from html error: failed to build %s: %v", fullPath, err)
	}

	composite := render.CompositeResult{}
	for _, node := range renderTree {
		res, err := node.Render(ctx)
		if err != nil {
			log.Fatal(err)
		}
		composite.Children = append(composite.Children, res)
	}

	return composite, nil
}

func (t *TemplateNode) AddChildren(node render.Node) {
	panic("Template html tag must be self closed and can't have children")
}
