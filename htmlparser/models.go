package htmlparser

type HtmlParser interface {
	ParseHtml(html string) ([]*HtmlTag, error)
}

type HtmlTag struct {
	Name          string
	InnerHtml     string
	InnerContent  string
	IsSelfClosing bool
	Attributes    []HtmlAttribute
	Parent        *HtmlTag
	Children      []*HtmlTag
	Pos           Position
	htmlStart     int
}

type HtmlAttribute struct {
	Name         string
	Value        string
	IsValueExist bool
}

type Position struct {
	Line   int
	Column int
}
