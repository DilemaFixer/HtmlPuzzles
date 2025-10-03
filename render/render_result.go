package render

type RenderResult interface {
	ToNodes() (HtmlNodes, error)
}

type HtmlResult struct {
	Nodes HtmlNodes
}

func (r HtmlResult) ToNodes() (HtmlNodes, error) {
	return r.Nodes, nil
}

type CompositeResult struct {
	Children []RenderResult
}

func (c CompositeResult) ToNodes() (HtmlNodes, error) {
	all := make(HtmlNodes, 0)
	for _, ch := range c.Children {
		nodes, err := ch.ToNodes()
		if err != nil {
			return nil, err
		}
		all = append(all, nodes...)
	}
	return all, nil
}

type HostResult struct {
	Host     *HtmlNode
	Children CompositeResult
}

func (r HostResult) ToNodes() (HtmlNodes, error) {
	hCopy, err := r.Host.CloneDown(1)
	if err != nil {
		return nil, err
	}

	children, err := r.Children.ToNodes()
	if err != nil {
		return nil, err
	}
	hCopy.Children = children

	return HtmlNodes{hCopy}, nil
}

type AsyncResult struct {
	Future *Future[RenderResult]
}

func (r AsyncResult) ToNodes() (HtmlNodes, error) {
	results, err := r.Future.Get()
	if err != nil {
		return nil, err
	}
	nodes, err := results.ToNodes()
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
