package main

import (
	"fmt"

	htmlparser "github.com/DilemaFixer/HtmlParser"
	"github.com/DilemaFixer/HtmlPuzzles/render"
	"github.com/DilemaFixer/HtmlPuzzles/tools"
)

const (
	for_keyword        = "for"
	IterationsCountKey = "iterations_count"
)

type ForNode struct {
	count     uint64
	childrens []render.RenderNode
}

func NewForNode(count uint64) *ForNode {
	return &ForNode{
		count:     count,
		childrens: make([]render.RenderNode, 0),
	}
}

func (fNode *ForNode) Render(ctx *tools.Context) ([]*htmlparser.HtmlTag, error) {
	ctx.LayerUp()
	result := make([]*htmlparser.HtmlTag, 0, len(fNode.childrens))

	for i := 0; i < int(fNode.count); i++ {
		for _, children := range fNode.childrens {
			childrenHtml, err := children.Render(ctx)
			if err != nil {
				return nil, err
			}
			result = append(result, childrenHtml...)
		}
	}

	ctx.LayerDown()
	return result, nil
}

func (hNode *ForNode) AddChildren(node render.RenderNode) {
	if node == nil {
		return
	}
	hNode.childrens = append(hNode.childrens, node)
}

func forValidator(htmlNode *htmlparser.HtmlTag) error {
	if !htmlNode.HasAttribute(IterationsCountKey) {
		return fmt.Errorf("Target tag %s have'n attr %s", htmlNode.Name, IterationsCountKey)
	}
	return nil
}

func forParser(htmlNode *htmlparser.HtmlTag) (render.RenderNode, error) {
	attr := htmlNode.GetAttribute(IterationsCountKey)
	attrValue, err := attr.AsUint64()
	if err != nil {
		return nil, err
	}

	return NewForNode(attrValue), nil
}
