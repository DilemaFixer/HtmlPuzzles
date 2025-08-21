package htmlparser

import "fmt"

func (t *HtmlTag) CloneUp(maxDepth int, ignoreDepthLimit bool) (*HtmlTag, error) {
	if maxDepth < 0 {
		return nil, fmt.Errorf("Html tags clone up error: max depth is less than zero : %d", maxDepth)
	}

	root := t
	var targetDepth int = 0

	for {
		if targetDepth == maxDepth {
			break
		}

		if root.Parent == nil {
			if ignoreDepthLimit {
				break
			}
			return nil, fmt.Errorf("Html tags clone up error: maxDepth %d exceeds available parents", maxDepth)
		}

		root = root.Parent
		targetDepth++
	}

	if targetDepth == 0 {
		return t.clone(0, 0), nil
	}

	return root.clone(0, targetDepth), nil
}

func (t *HtmlTag) CloneDown(maxDepth int) (*HtmlTag, error) {
	if maxDepth < 1 {
		return nil, fmt.Errorf("Html tags clone down error: max depth is less than zero : %d", maxDepth)
	}

	return t.clone(0, maxDepth), nil
}

func (t *HtmlTag) clone(currentLayer, targetDepth int) *HtmlTag {
	if t == nil {
		return nil
	}

	clone := &HtmlTag{
		Name:          t.Name,
		InnerHtml:     t.InnerHtml,
		InnerContent:  t.InnerContent,
		IsSelfClosing: t.IsSelfClosing,
		Pos:           t.Pos,
		htmlStart:     t.htmlStart,
	}

	if len(t.Attributes) > 0 {
		clone.Attributes = make([]HtmlAttribute, len(t.Attributes))
		copy(clone.Attributes, t.Attributes)
	}

	if len(t.Children) > 0 && currentLayer < targetDepth {
		clone.Children = make([]*HtmlTag, len(t.Children))
		for i, child := range t.Children {
			childClone := child.clone(currentLayer+1, targetDepth)
			if childClone != nil {
				childClone.Parent = clone
			}
			clone.Children[i] = childClone
		}
	}

	return clone
}
