package html

import "testing"

func TestParser(t *testing.T) {
	html := `<!--meta:role=button-->
<div>hi</div>
`
	tree, err := ParseHtml(html)
	if err != nil {
		t.Fatal(err)
	}
	for _, br := range tree {
		PrintTree(br)
	}
}
