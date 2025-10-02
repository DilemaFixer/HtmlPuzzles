package render

type NodeParser interface {
	GetTarget() string
	Parser(htmlNode *HtmlNode, childrenCount uint64) (Node, error)
}
