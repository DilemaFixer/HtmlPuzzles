package tagrender

import (
	"fmt"

	"github.com/DilemaFixer/HtmlPuzzles/consts"
	"github.com/DilemaFixer/HtmlPuzzles/htmlparser"
)

func IsRendererTag(name string) bool {
	switch name {
	case "for":
		return true
	case "set":
		return true
	}
	return false
}

func BuildRenderer(node *htmlparser.HtmlTag) (Renderer, error) {

	if node.Name == "for" {
		if node.HasAttribute(consts.IterationsCountKey) {
			attr := node.GetAttribute(consts.IterationsCountKey)
			iterCount, err := attr.AsUint64()
			if err != nil {
				return nil, fmt.Errorf("Error html tag attrebute type custing: %s", err.Error())
			}

			iterator, err := NewConstantIterator(iterCount)
			if err != nil {
				return nil, fmt.Errorf("Error building constant iterator for 'for' renderer : %s", err.Error())
			}

			return NewForRenderer(iterator), nil
		}
		return nil, fmt.Errorf("Renderer building error: Exist target attribute %s", consts.IterationsCountKey)
	}

	if node.Name == "set" {
		if node.HasAttribute(consts.WayKey) {
			attr := node.GetAttribute(consts.WayKey)
			seter := NewSetRenderer(attr.AsString())
			return seter, nil
		}

		return nil, fmt.Errorf("Renderer building error: Exist target attribute %s", consts.WayKey)
	}

	return nil, fmt.Errorf("Renderer building error: Undefinde renderer type %s", node.Name)
}
