package render

import (
	"fmt"
)

type TreeBuilder struct {
	binds          map[string]NodeParser
	defaultParser  NodeParser
	rebindsAllowed bool
}

func NewTreeBuilder(defaultParser NodeParser, rebinsAllowed bool) *TreeBuilder {
	return &TreeBuilder{
		binds:          make(map[string]NodeParser),
		rebindsAllowed: rebinsAllowed,
		defaultParser:  defaultParser,
	}
}

func (builder *TreeBuilder) Bind(parser NodeParser) error {
	_, isExist := builder.binds[parser.GetTarget()]
	if isExist && !builder.rebindsAllowed {
		return fmt.Errorf("render tree builder error: rebinds not allowed , try rebind %s", parser.GetTarget())
	}

	builder.binds[parser.GetTarget()] = parser
	return nil
}

func (builder *TreeBuilder) HasBind(target string) bool {
	_, isExist := builder.binds[target]
	return isExist
}

func (builder *TreeBuilder) Build(source HtmlNodes) ([]Node, error) {
	result := make([]Node, 0)
	for _, branch := range source {
		node, err := builder.BuildBranch(branch)
		if err != nil {
			return nil, err
		}
		result = append(result, node)
	}
	return result, nil
}

func (builder *TreeBuilder) BuildBranch(sourceBranch *HtmlNode) (Node, error) {
	var (
		node Node
		err  error
	)

	if builder.HasBind(sourceBranch.Name) {
		node, err = builder.parseNodeFromBind(sourceBranch)
	} else {
		node, err = builder.parseNodeFromDefault(sourceBranch)
	}
	if err != nil {
		return nil, err
	}

	if err := builder.parseChildren(node, sourceBranch); err != nil {
		return nil, err
	}

	return node, nil
}

func (builder *TreeBuilder) parseNodeFromBind(source *HtmlNode) (Node, error) {
	parser := builder.binds[source.Name]
	node, err := parser.Parser(source, uint64(len(source.Children)))
	if err != nil {
		return nil, fmt.Errorf("render branch building error: parse branch %s failed with error \n%v", source.Name, err)
	}
	return node, nil
}

func (builder *TreeBuilder) parseNodeFromDefault(source *HtmlNode) (Node, error) {
	node, err := builder.defaultParser.Parser(source, uint64(len(source.Children)))
	if err != nil {
		return nil, fmt.Errorf("render branch building error: parse branch '%s' failed with error \n'%v'", source.Name, err)
	}
	return node, nil
}

func (builder *TreeBuilder) parseChildren(parent Node, sourceBranch *HtmlNode) error {
	for _, subBranche := range sourceBranch.Children {
		subNode, err := builder.BuildBranch(subBranche)
		if err != nil {
			return err
		}
		parent.AddChildren(subNode)
	}
	return nil
}
