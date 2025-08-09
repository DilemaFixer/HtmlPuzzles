package main

import (
	html "github.com/DilemaFixer/HtmlPuzzles/htmlparser"
)

func main() {
	root, err := html.ParseHtml(`
<div class="container">
    <h1>Hello World</h1>
    <p id="text">Some content</p>
    <img src="image.jpg" alt="Image"/>
</div>
`)

	if err == nil {
		html.PrintHtmlTree(root)
	}
}
