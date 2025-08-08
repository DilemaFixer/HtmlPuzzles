package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	htmlContent := `
        <container>
            <header>My Site</header>
            <for items=".Users" item="user">
                <card>
                    <h3><set>user.Name</set></h3>
                    <p><set>user.Email</set></p>
                </card>
            </for>
            <if condition=".User.IsAdmin">
                <admin-panel>
                    <button>Delete User</button>
                </admin-panel>
            </if>
        </container>
    `

	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		panic(err)
	}

	printAllTags(doc, 0)
}

func printAllTags(n *html.Node, depth int) {
	if n == nil {
		return
	}

	indent := strings.Repeat("  ", depth)

	switch n.Type {
	case html.ElementNode:
		fmt.Printf("%s<%s>\n", indent, n.Data)

		// Получаем HTML содержимое этого элемента (innerHTML)
		innerHTML := getInnerHTML(n)
		if innerHTML != "" {
			fmt.Printf("%s  InnerHTML: %s\n", indent, innerHTML)
		}

	case html.TextNode:
		text := strings.TrimSpace(n.Data)
		if text != "" {
			fmt.Printf("%sText: %q\n", indent, text)
		}
	}

	// Рекурсивно обходим дочерние элементы
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		printAllTags(c, depth+1)
	}
}

func getInnerHTML(n *html.Node) string {
	var result strings.Builder

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		renderNode(c, &result)
	}

	return strings.TrimSpace(result.String())
}

// renderNode рендерит узел в HTML строку
func renderNode(n *html.Node, w *strings.Builder) {
	switch n.Type {
	case html.TextNode:
		w.WriteString(n.Data)

	case html.ElementNode:
		w.WriteString("<")
		w.WriteString(n.Data)

		// Добавляем атрибуты
		for _, attr := range n.Attr {
			w.WriteString(" ")
			w.WriteString(attr.Key)
			w.WriteString("=\"")
			w.WriteString(attr.Val)
			w.WriteString("\"")
		}

		// Проверяем, есть ли дочерние элементы
		if n.FirstChild != nil {
			w.WriteString(">")

			// Рендерим дочерние элементы
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				renderNode(c, w)
			}

			w.WriteString("</")
			w.WriteString(n.Data)
			w.WriteString(">")
		} else {
			// Самозакрывающийся тег
			w.WriteString("/>")
		}

	case html.CommentNode:
		w.WriteString("<!--")
		w.WriteString(n.Data)
		w.WriteString("-->")

	case html.DoctypeNode:
		w.WriteString("<!DOCTYPE ")
		w.WriteString(n.Data)
		w.WriteString(">")
	}
}

// getOuterHTML возвращает HTML включая сам элемент (outerHTML)
func getOuterHTML(n *html.Node) string {
	var result strings.Builder
	renderNode(n, &result)
	return result.String()
}

// Пример использования для поиска конкретного элемента
func findElementByTag(n *html.Node, tagName string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tagName {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := findElementByTag(c, tagName); found != nil {
			return found
		}
	}

	return nil
}
