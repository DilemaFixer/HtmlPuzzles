package html

import "fmt"

type HtmlTag struct {
	Name          string
	InnerHtml     string
	InnerContent  string
	IsSelfClosing bool
	Attributes    map[string]HtmlAttribute
	Parent        *HtmlTag
	Children      []*HtmlTag
	Pos           Position

	Meta      map[string]string
	htmlStart int
}

type Position struct {
	Line   int
	Column int
}

func (t *HtmlTag) CloneUp(maxDepth int, ignoreDepthLimit bool) (*HtmlTag, error) {
	if maxDepth < 0 {
		return nil, fmt.Errorf("Html tags clone up error: max depth is less than zero : %d", maxDepth)
	}

	root := t
	depth := 0

	for depth < maxDepth {
		if root.Parent == nil {
			if ignoreDepthLimit {
				break
			}
			return nil, fmt.Errorf("Html tags clone up error: maxDepth %d exceeds available parents", maxDepth)
		}
		root = root.Parent
		depth++
	}

	if depth == 0 {
		return t.clone(0, 0), nil
	}
	return root.clone(0, depth), nil
}

func (t *HtmlTag) CloneDown(maxDepth int) (*HtmlTag, error) {
	if maxDepth < 1 {
		return nil, fmt.Errorf("Html tags clone down error: max depth is less than one : %d", maxDepth)
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
		clone.Attributes = make(map[string]HtmlAttribute, len(t.Attributes))
		for key, value := range t.Attributes {
			clone.Attributes[key] = value
		}
	}

	if len(t.Children) > 0 && currentLayer < targetDepth {
		clone.Children = make([]*HtmlTag, 0, len(t.Children))
		for _, child := range t.Children {
			if childClone := child.clone(currentLayer+1, targetDepth); childClone != nil {
				childClone.Parent = clone
				clone.Children = append(clone.Children, childClone)
			}
		}
	}

	return clone
}

func (t *HtmlTag) RemoveChild(child *HtmlTag) {
	if t == nil || child == nil || len(t.Children) == 0 {
		return
	}

	for i, c := range t.Children {
		if c == child {
			c.Parent = nil
			t.Children = append(t.Children[:i], t.Children[i+1:]...)
			return
		}
	}
}

func (t *HtmlTag) HasAttribute(name string) bool {
	_, exists := t.Attributes[name]
	return exists
}

func (t *HtmlTag) GetAttribute(name string) *HtmlAttribute {
	if attr, exists := t.Attributes[name]; exists {
		return &attr
	}
	return nil
}

func (t *HtmlTag) SetAttribute(name, value string) {
	if t.Attributes == nil {
		t.Attributes = make(map[string]HtmlAttribute)
	}
	attr := HtmlAttribute{
		Name:         name,
		Value:        value,
		IsValueExist: true,
	}
	t.Attributes[name] = attr
}

func (t *HtmlTag) RemoveAttribute(name string) {
	if t.Attributes == nil {
		return
	}
	delete(t.Attributes, name)
}
