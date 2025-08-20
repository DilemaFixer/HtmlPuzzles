package main

import "github.com/DilemaFixer/HtmlPuzzles/htmlrender"

type One struct {
	Dummy int
	T     Two
}

type Two struct {
	Dummy string
	T     Three
}

type Three struct {
	Dummy bool
	F     Four
}

type Four struct {
	Dummy float64
	Four  int
}

var HTML string = `<div> <for iterations_count=10> <h1> <set way="T.T.F.Four"/> </h1> </for> </div>`

func main() {
	/*
		 one := One{
				T: Two{
					T: Three{
						F: Four{
							Four: 42,
						},
					},
				},
			}
			path := "T.T.F.Four"

			ctx := render.NewContext(&one)
			seter := tagrender.NewSetRenderer(path)

			result, err := seter.Render(ctx)
			if err != nil {
				fmt.Println("Error rendering:", err)
				return
			}
			fmt.Println("Rendered result:", result)
	*/

	htmlrender.ParseToRendererNodes(HTML)
}
