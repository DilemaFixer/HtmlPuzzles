package html

import (
	"fmt"
	"sort"
	"strings"
)

func PrintTree(tag *HtmlTag) {
	printHtmlTreeRecursive(tag, 0)
}

func PrintTreeRecursive(tag *HtmlTag, depth int) {
	printHtmlTreeRecursive(tag, depth)
}

func printHtmlTreeRecursive(tag *HtmlTag, depth int) {
	if tag == nil {
		return
	}

	indent := strings.Repeat("  ", depth)
	fmt.Printf("%s<%s", indent, tag.Name)

	keys := make([]string, 0, len(tag.Attributes))
	for k := range tag.Attributes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		attr := tag.Attributes[k]
		if attr.IsValueExist {
			fmt.Printf(" %s=\"%s\"", attr.Name, attr.Value)
		} else {
			fmt.Printf(" %s", attr.Name)
		}
	}

	if tag.IsSelfClosing {
		fmt.Printf("/>\n")
		return
	}

	fmt.Printf(">")

	if tag.InnerContent != "" {
		fmt.Printf("%s", tag.InnerContent)
	}

	if len(tag.Children) > 0 {
		fmt.Printf("\n")
		for _, child := range tag.Children {
			printHtmlTreeRecursive(child, depth+1)
		}
		fmt.Printf("%s", indent)
	}

	fmt.Printf("</%s>\n", tag.Name)
}
