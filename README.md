> [!WARNING]
> The library is not in contributions ready state yet , test not exist

# HtmlPuzzles

**HtmlPuzzles** is an experimental Go library for creating custom HTML components with a custom render pipeline. This project was created for educational purposes to study HTML parsing, DSL creation, and architectural patterns.

## 🎯 Main Idea

The library allows you to:
- Create custom HTML tags with custom logic
- Implement your own components (loops, conditions, variables, etc.)
- Manage execution context with layer support
- Parse and render HTML with custom tags

## 📦 Installation

```bash
go get github.com/DilemaFixer/HtmlPuzzles
go get github.com/DilemaFixer/HtmlParser
```

## 🏗️ Architecture

### Core Components:

- **TagsRenderer** - main renderer, manages parsing and rendering
- **RenderNode** - interface for all render nodes
- **HostNode** - wrapper for regular HTML tags
- **Context** - execution context with layer support and typed data

### Workflow:

1. HTML is parsed into AST using HtmlParser
2. AST is converted into a RenderNode tree
3. Each node is rendered with context passing
4. Result is serialized back to HTML

## 🚀 Quick Start

### Simple example with a loop:

```go
package main

import (
    "fmt"
    "github.com/DilemaFixer/HtmlPuzzles/render"
)

const HTML = `
<div>
    <for iterations_count=3>
        <h1>Hello World!</h1>
    </for>
</div>`

func main() {
    renderer := render.NewTagsRenderer()
    renderer.Bind("for", forParser, forValidator)
    
    result, err := renderer.RenderHtml(HTML)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(result)
    // Output: <div><h1>Hello World!</h1><h1>Hello World!</h1><h1>Hello World!</h1></div>
}
```

## 🔧 Creating Custom Components

### Step 1: Implement RenderNode

```go
type CustomNode struct {
    // your fields
    childrens []render.RenderNode
}

func (n *CustomNode) Render(ctx *tools.Context) ([]*htmlparser.HtmlTag, error) {
    ctx.LayerUp() // move to new context layer
    defer ctx.LayerDown() // return to previous layer
    
    // your rendering logic
    result := make([]*htmlparser.HtmlTag, 0)
    
    // render child elements
    for _, child := range n.childrens {
        childHtml, err := child.Render(ctx)
        if err != nil {
            return nil, err
        }
        result = append(result, childHtml...)
    }
    
    return result, nil
}

func (n *CustomNode) AddChildren(node render.RenderNode) {
    if node != nil {
        n.childrens = append(n.childrens, node)
    }
}
```

### Step 2: Validator

```go
func customValidator(htmlNode *htmlparser.HtmlTag) error {
    // check required attributes
    if !htmlNode.HasAttribute("required_attr") {
        return fmt.Errorf("tag %s must contain required_attr attribute", htmlNode.Name)
    }
    return nil
}
```

### Step 3: Parser

```go
func customParser(htmlNode *htmlparser.HtmlTag) (render.RenderNode, error) {
    // extract data from attributes
    attr := htmlNode.GetAttribute("required_attr")
    value, err := attr.AsString() // or AsInt64(), AsUint64(), AsBool()
    if err != nil {
        return nil, err
    }
    
    return &CustomNode{
        // initialize with extracted data
        childrens: make([]render.RenderNode, 0),
    }, nil
}
```

### Step 4: Registration

```go
renderer := render.NewTagsRenderer()
renderer.Bind("custom", customParser, customValidator)
```

## 📋 Complete Example: Loop Component

```go
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

// Node for loop
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
    defer ctx.LayerDown()
    
    result := make([]*htmlparser.HtmlTag, 0)

    // execute loop
    for i := 0; i < int(fNode.count); i++ {
        // can add iteration variable to context
        ctx.SetIntLayered("index", int64(i))
        
        for _, children := range fNode.childrens {
            childrenHtml, err := children.Render(ctx)
            if err != nil {
                return nil, err
            }
            result = append(result, childrenHtml...)
        }
    }

    return result, nil
}

func (fNode *ForNode) AddChildren(node render.RenderNode) {
    if node != nil {
        fNode.childrens = append(fNode.childrens, node)
    }
}

// Validator
func forValidator(htmlNode *htmlparser.HtmlTag) error {
    if !htmlNode.HasAttribute(IterationsCountKey) {
        return fmt.Errorf("tag %s must contain %s attribute", 
            htmlNode.Name, IterationsCountKey)
    }
    return nil
}

// Parser
func forParser(htmlNode *htmlparser.HtmlTag) (render.RenderNode, error) {
    attr := htmlNode.GetAttribute(IterationsCountKey)
    attrValue, err := attr.AsUint64()
    if err != nil {
        return nil, err
    }

    return NewForNode(attrValue), nil
}

// Usage
func main() {
    const HTML = `
    <div>
        <for iterations_count=3>
            <p>Iteration</p>
        </for>
    </div>`

    renderer := render.NewTagsRenderer()
    renderer.Bind(for_keyword, forParser, forValidator)
    
    htmlAsString, err := renderer.RenderHtml(HTML)
    if err != nil {
        fmt.Printf("Error: %s", err.Error())
        return
    }
    
    fmt.Println(htmlAsString)
}
```

## 🗃️ Working with Context

Context supports layered architecture and typed data:

```go
ctx := tools.NewContext()

// Setting values
ctx.SetString("name", "John")
ctx.SetInt("age", 25)
ctx.SetBool("active", true)

// Working with layers
ctx.LayerUp()
ctx.SetStringLayered("localVar", "value") // will be available only on this layer
ctx.LayerDown()

// Getting values
name := ctx.String("name")
age := ctx.Int("age")
localVar := ctx.StringLayered("localVar") // empty string since layer changed

// Checking existence
if ctx.HasString("name") {
    // value exists
}
```

### Supported types:
- `string`
- `int64` 
- `uint64`
- `float64`
- `bool`
- `any` (objects)

## 💡 Component Examples

### Conditional Rendering

```go
type IfNode struct {
    condition string
    childrens []render.RenderNode
}

func (n *IfNode) Render(ctx *tools.Context) ([]*htmlparser.HtmlTag, error) {
    // check condition from context
    if ctx.Bool(n.condition) {
        result := make([]*htmlparser.HtmlTag, 0)
        for _, child := range n.childrens {
            childHtml, err := child.Render(ctx)
            if err != nil {
                return nil, err
            }
            result = append(result, childHtml...)
        }
        return result, nil
    }
    return make([]*htmlparser.HtmlTag, 0), nil
}
```

### Variable Component

```go
type SetNode struct {
    key   string
    value string
    childrens []render.RenderNode
}

func (n *SetNode) Render(ctx *tools.Context) ([]*htmlparser.HtmlTag, error) {
    ctx.LayerUp()
    defer ctx.LayerDown()
    
    // set variable in context
    ctx.SetStringLayered(n.key, n.value)
    
    // render child elements
    result := make([]*htmlparser.HtmlTag, 0)
    for _, child := range n.childrens {
        childHtml, err := child.Render(ctx)
        if err != nil {
            return nil, err
        }
        result = append(result, childHtml...)
    }
    return result, nil
}
```

## ⚠️ Limitations

- Project created for educational purposes
- Not intended for production use
- Basic functionality, requires refinement for complex scenarios
- Missing handling of many edge cases

## 🤝 TODO

- [ ] Add more secure layer key generation
- [ ] Improve error handling
- [ ] Add more component examples
