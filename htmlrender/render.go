package htmlrender

import (
	"github.com/DilemaFixer/HtmlPuzzles/htmlparser"
	"github.com/DilemaFixer/HtmlPuzzles/htmlrender/tagrender"
	"github.com/DilemaFixer/HtmlPuzzles/tools"
)

type TesetRenderer struct{}

func (t *TesetRenderer) Render(ctx tools.Context) (*htmlparser.HtmlTag, error) {
	return nil, nil
}

func (t *TesetRenderer) AddChildren(*tagrender.Renderer) {
}

func ParseToRendererNodes(html string) ([]tagrender.Renderer, error) {
	htmlNodes, err := htmlparser.ParseHtml(html)

	if err != nil {
		return nil, err
	}

	result := make([]tagrender.Renderer, 0)

	for _, node := range htmlNodes {
		renderer, err := handlingHtmlTree(node, 0)
		if err != nil {
			return nil, err
		}
		result = append(result, renderer)
	}

	return result, nil
}

func handlingHtmlTree(node *htmlparser.HtmlTag, deep uint) (tagrender.Renderer, error) {
	deep++
	if tagrender.IsRendererTag(node.Name) {
		renderer, err := tagrender.BuildRenderer(node)
		if err != nil {
			return nil, err
		}

		if deep != 0 {
			var targetPerent *htmlparser.HtmlTag = node

			var i uint
			for i = 0; i < deep-1; i++ {
				targetPerent = targetPerent.Parent
			}

			copy := targetPerent.DeepClone(deep - 1)
			renderer.EstablishResponsibilityForRendering(copy)
			deep = 0
		}

		if node.Children != nil {
			for _, children := range node.Children {
				childrenRenderer, err := handlingHtmlTree(children, deep)
				if err != nil {
					return nil, err
				}
				renderer.AddChildren(childrenRenderer)
			}
		}

		return renderer, nil
	}

	renderer, err := handlingHtmlTree(node, deep)
	if err != nil {
		return nil, err
	}
	deep--
	return renderer, nil
}
