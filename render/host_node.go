package render

import (
	htmlparser "github.com/DilemaFixer/HtmlParser"
	"github.com/DilemaFixer/HtmlPuzzles/tools"
)

type HostNode struct {
	htmlNode  *htmlparser.HtmlTag
	childrens []RenderNode
}

func NewHostNode(htmlNode *htmlparser.HtmlTag) *HostNode {
	return &HostNode{
		htmlNode:  htmlNode,
		childrens: make([]RenderNode, 0),
	}
}

func (hNode *HostNode) Render(ctx *tools.Context) ([]*htmlparser.HtmlTag, error) {
	ctx.LayerUp()
	nodeCopy, err := hNode.htmlNode.CloneUp(0, false)

	if err != nil {
		return nil, err
	}

	for _, children := range hNode.childrens {
		childrensHtml, err := children.Render(ctx)
		if err != nil {
			return nil, err
		}
		nodeCopy.Children = append(nodeCopy.Children, childrensHtml...)
	}

	ctx.LayerDown()
	result := make([]*htmlparser.HtmlTag, 0, 1)
	result = append(result, nodeCopy)
	return result, nil
}

func (hNode *HostNode) AddChildren(node RenderNode) {
	if node == nil {
		return
	}
	hNode.childrens = append(hNode.childrens, node)
}
