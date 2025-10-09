> [!WARNING]
> The library is not in contributions ready state yet , test not exist

# HtmlPuzzles

**HtmlPuzzles** is an experimental Go library for creating custom HTML components with a custom render pipeline. This project was created for educational purposes to study HTML parsing, DSL creation, and architectural patterns.

##  Main Idea

The library allows you to:
- Create custom HTML tags with custom logic
- Implement your own components (loops, conditions, variables, etc.)
- Manage execution context with layer support
- Parse and render HTML with custom tags

##  Installation

```bash
go get github.com/DilemaFixer/HtmlPuzzles
```

For more details chack [wiki page](https://github.com/DilemaFixer/HtmlPuzzles/wiki)

## Example of possibilities

Input html : 
``` Html
<!-- <sync> - start rendering inner html in separate gorutine without stoping main -->
<!-- <for> - just repiting inner html , requared 'itr_count' attrebute  -->
<!-- <templ> - render .html file as innner content -->
<!-- <wrapper> - using as templ for html tag -->
<!-- 'wrapped' attrebute requared always -->
<!-- all attrebutes will be move to tag that describe in wrapped -->
<!-- if attrebute start with ':' prefix, value will be geting from runtime context , you can describe way to prop -->
<!-- :src="user.AvatarUrl|string" -> src="path/to/img/from/prop" -->
<div>
    <sync>
        <for itr_count=3>
            <h1>Hello World!</h1>
        </for>
    </sync>
    <sync>
        <templ source="index.html" />
    </sync>
    <wrapper wrapped="img" :src="user.AvatarUrl|string" :wight="user.Wight|uint" />
</div>
```

Template file index.html :
``` Html
<div>
    <for itr_count=5>
        <h1>
            Hello from template!!!
        </h1>
    </for>
</div>
```

Result :
``` Html
<div>
   <h1>Hello World!</h1>
   <h1>Hello World!</h1>
   <h1>Hello World!</h1>
   <div>
      <h1>Hello from template!!!</h1>
      <h1>Hello from template!!!</h1>
      <h1>Hello from template!!!</h1>
      <h1>Hello from template!!!</h1>
      <h1>Hello from template!!!</h1>
   </div>
   <img src="test/img/bitch" wight="1"/>
</div>
```
#How create your own HTML tag

## Parser

First of all, you need to understand how the logic is distributed. In order for the library to recognize that your tag exists, you need to create a parser that will implement the following interface: 

``` Go
type NodeParser interface {
	GetTarget() string
	Parser(htmlNode *HtmlNode, childrenCount uint64) (Node, error)
}
```

To make it clearer, I will show you an example of implementing the for tag : 

``` Go
import (
	"fmt"

	"github.com/DilemaFixer/HtmlPuzzles/examples/nodes"
	"github.com/DilemaFixer/HtmlPuzzles/render"
)

type ForParser struct{}

func NewForParser() render.NodeParser {
	return &ForParser{}
}

func (f *ForParser) GetTarget() string {
	return "for"
}

func (f *ForParser) Parser(htmlNode *render.HtmlNode, childrenCount uint64) (render.Node, error) {
	if err := validateNode(htmlNode); err != nil {
		return nil, err
	}

	attr := htmlNode.GetAttribute("itr_count")
	itr_count, err := attr.AsUint64()
	if err != nil {
		return nil, err
	}
	return nodes.NewForNode(itr_count, childrenCount), nil
}

func validateNode(htmlNode *render.HtmlNode) error {
	if !htmlNode.HasAttribute("itr_count") {
		return fmt.Errorf("parsing 'for' html tag error: expected attribute 'itr_count' but it not exists")
	}
	return nil
}
```

In the GetTarget method, you constantly write the name of your tag (it is good practice to put the name in a private constant) :

``` Go
func (f *ForParser) GetTarget() string {
	return "for"
}
```

Next, you implement the Parser method, which accepts an HTML tag and the expected number of children (optimization). The HTML node object itself has enough methods to allow you to get the data you need to build the render.Node structure : 

``` Go
func (f *ForParser) Parser(htmlNode *render.HtmlNode, childrenCount uint64) (render.Node, error) {
	if err := validateNode(htmlNode); err != nil {
		return nil, err
	}

	attr := htmlNode.GetAttribute("itr_count")
	itr_count, err := attr.AsUint64()
	if err != nil {
		return nil, err
	}
	return nodes.NewForNode(itr_count, childrenCount), nil
}
```

I would like to highlight the good practice of moving HTML tag validation to a separate validateNode method that returns an error, which makes it syntactically convenient to call and handle the error : 

``` Go
func validateNode(htmlNode *render.HtmlNode) error {
	if !htmlNode.HasAttribute("itr_count") {
		return fmt.Errorf("parsing 'for' html tag error: expected attribute 'itr_count' but it not exists")
	}
	return nil
}
```

using like :

``` Go
if err := validateNode(htmlNode); err != nil {
		return nil, err
}
```


## Node 

Node interface : 

``` Go
type Node interface {
	Render(ctx *Context) (RenderResult, error)
	AddChildren(node Node)
}
```

The same tag for, but now implemented by the node responsible for processing it : 

``` Go
type ForNode struct {
	children  []render.Node
	itr_count uint64
}

func NewForNode(itr_count uint64, childrenCount uint64) render.Node {
	return &ForNode{
		itr_count: itr_count,
		children:  make([]render.Node, 0, childrenCount),
	}
}

func (f *ForNode) Render(ctx *render.Context) (render.RenderResult, error) {
	ctx.LayerUp()
	defer ctx.LayerDown()

	result := render.CompositeResult{}
	for i := uint64(0); i < f.itr_count; i++ {
		for _, subNode := range f.children {
			rendered, err := subNode.Render(ctx)
			if err != nil {
				return nil, err
			}
			result.Children = append(result.Children, rendered)
		}
	}

	return result, nil
}

func (f *ForNode) AddChildren(node render.Node) {
	if node == nil {
		return
	}
	f.children = append(f.children, node)
}

```

First, implement the structure and constructor. The structure must contain the `children[]render.Node` field, and the constructor must accept `childrenCount uint64` and use this value to initialize the array of child elements children: `make([]render.Node, 0, childrenCount)`

``` Go
type ForNode struct {
	children  []render.Node
	itr_count uint64
}

func NewForNode(itr_count uint64, childrenCount uint64) render.Node {
	return &ForNode{
		itr_count: itr_count,
		children:  make([]render.Node, 0, childrenCount),
	}
}
```

I think the implementation of `AddChildren(node render.Node)` needs no explanation :

``` Go
func (f *ForNode) AddChildren(node render.Node) {
	if node == nil {
		return
	}
	f.children = append(f.children, node)
}
```

## Render Result

To continue the explanation, I need to talk about RenderResult. This abstraction was created to solve the problem of asynchronous processing of a separate branch. Initially, the Render method simply returned an array of children that it processes, i.e., it required an immediate response from the render node, but this abstraction allowed processing to be postponed and also transferred the formatting and tree assembly functionality to separate structures. Next, we will show the interface and several implementations of RenderResult: 

``` Go
type RenderResult interface {
	ToNodes() (HtmlNodes, error)
}
```

``` Go
type HtmlResult struct {
	Nodes HtmlNodes
}

func (r HtmlResult) ToNodes() (HtmlNodes, error) {
	return r.Nodes, nil
}
```

``` Go
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
```

``` Go
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
```

``` Go
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
```

Next, we can talk about the `Render` method. This is where you will implement the core logic for creating, processing, and modifying the HTML tree. I can't give you a clear template for what to do here; you are free to choose! But to better understand the possibilities, you should take a look at the examples folder in the project root. 

