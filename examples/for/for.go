package main

import (
	"fmt"

	"github.com/DilemaFixer/HtmlPuzzles/render"
)

const HTML = `  <div>
					<for iterations_count=10>
						<h1>
							<set way="T.T.Data" wrapper="h1"/>
						</h1>
					</for>
				</div>`

func main() {
	renderer := render.NewTagsRenderer()
	renderer.Bind(for_keyword, forParser, forValidator)
	htmlAsString, err := renderer.RenderHtml(HTML)
	if err != nil {
		fmt.Printf("Err: %s", err.Error())
	}
	fmt.Println(htmlAsString)
}
