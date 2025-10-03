package main

import (
	"fmt"
	"log"

	"github.com/DilemaFixer/HtmlPuzzles/examples/parsers"
	"github.com/DilemaFixer/HtmlPuzzles/html"
	"github.com/DilemaFixer/HtmlPuzzles/render"
)

const (
	TemplateRootDir = "./examples/templates"
	UserObjectKey   = "user"
)

const InputHTML = `
<div>
	<sync>
		<for itr_count=3>
			<h1>Hello World!</h1>
		</for>
	</sync>
	<sync>
		<templ source="index.html"/> 
	</sync>
	<wrapper wrapped="img" :src="user.AvatarUrl|string" :wight="user.Wight|uint"/>
</div>
`

type User struct {
	AvatarUrl string
	Wight     uint64
}

func main() {
	htmlAST, err := html.ParseHtml(InputHTML)
	if err != nil {
		log.Fatalf("failed to parse HTML: %v", err)
	}

	builder := render.NewTreeBuilder(parsers.NewHostParser(), false)
	if err := bindParsers(builder); err != nil {
		log.Fatalf("failed to bind parsers: %v", err)
	}

	renderedTree, err := builder.Build(htmlAST)
	if err != nil {
		log.Fatalf("failed to build render tree: %v", err)
	}

	ctx := render.NewContext(builder, TemplateRootDir)

	user := User{
		AvatarUrl: "test/img/bitch",
		Wight:     1,
	}
	ctx.SetObject(UserObjectKey, &user)

	result := renderTree(ctx, renderedTree)
	htmlOutput := html.SerializeToString(result)

	fmt.Println(htmlOutput)
}

func bindParsers(builder *render.TreeBuilder) error {
	parsersToBind := []render.NodeParser{
		parsers.NewForParser(),
		parsers.NewSyncParser(),
		parsers.NewTemplateParser(),
		parsers.NewWrapperParser(),
	}

	for _, p := range parsersToBind {
		if err := builder.Bind(p); err != nil {
			return err
		}
	}
	return nil
}

func renderTree(ctx *render.Context, nodes []render.Node) render.HtmlNodes {
	result := render.CompositeResult{}

	for _, node := range nodes {
		rendered, err := node.Render(ctx)
		if err != nil {
			log.Fatalf("failed to render node: %v", err)
		}
		result.Children = append(result.Children, rendered)
	}

	ctx.WaitAll()

	finalNodes, err := result.ToNodes()
	if err != nil {
		log.Fatalf("failed to convert composite result to nodes: %v", err)
	}
	return finalNodes
}
