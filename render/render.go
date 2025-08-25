package render

import (
	"fmt"

	htmltool "github.com/DilemaFixer/HtmlParser"
	"github.com/DilemaFixer/HtmlPuzzles/tools"
)

func (r *TagsRenderer) RenderHtml(html string) (string, error) {
	htmlParser := htmltool.NewHtmlParser()
	htmlAsts, err := htmlParser.ParseHtml(html)
	if err != nil {
		return "", fmt.Errorf("Rendering html error: %s", err.Error())
	}

	htmlRenderingBranches, err := r.HtmlToRenderTree(htmlAsts)
	if err != nil {
		return "", fmt.Errorf("Rendering html error: %s", err.Error())
	}
	ctx := tools.NewContext()
	renderedHtmlBranches := make([]*htmltool.HtmlTag, 0, len(htmlAsts))

	for _, brach := range htmlRenderingBranches {
		renderedHtmlBranch, err := brach.Render(ctx)
		if err != nil {
			return "", fmt.Errorf("Rendering html error: %s", err.Error())
		}
		renderedHtmlBranches = append(renderedHtmlBranches, renderedHtmlBranch...)
	}

	serializer := htmltool.NewHtmlSerializer()
	htmlAsString := serializer.RenderHtml(renderedHtmlBranches)
	return htmlAsString, nil
}
