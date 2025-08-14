package htmlparser

import (
	"fmt"
	"strings"

	"github.com/DilemaFixer/HtmlPuzzles/tools"
)

type HtmlParser struct {
	customHandlers map[string]func(*tools.Scanner) (*HtmlTag, error)
}

type HtmlTag struct {
	Name          string
	InnerHtml     string
	InnerContent  string
	IsSelfClosing bool
	Attributes    []HtmlAttribute
	Parent        *HtmlTag
	Children      []*HtmlTag
	Pos           Position
	htmlStart     int
}

type HtmlAttribute struct {
	Name         string
	Value        string
	IsValueExist bool
}

type Position struct {
	Line   int
	Column int
}

func NewHtmlParser() *HtmlParser {
	return &HtmlParser{
		customHandlers: make(map[string]func(*tools.Scanner) (*HtmlTag, error)),
	}
}

func (parser *HtmlParser) AddCustomAttributeHandler(tagName string, handler func(*tools.Scanner) (*HtmlTag, error)) {
	if handler == nil || strings.TrimSpace(tagName) == "" {
		return
	}

	parser.customHandlers[tagName] = handler
}

func ParseHtml(html string) ([]*HtmlTag, error) {
	parser := NewHtmlParser()
	return parser.ParseHtml(html)
}

func (parser *HtmlParser) ParseHtml(html string) ([]*HtmlTag, error) {
	scanner := tools.NewScanner(html)
	closeStack := tools.NewStack[string]()
	var roots []*HtmlTag
	var current *HtmlTag = nil

	for !scanner.EOF() {
		scanner.SkipWhitespace()
		if scanner.EOF() {
			break
		}

		if scanner.Ch == '<' {
			if scanner.PeekNext() == '/' {
				if closeStack.IsEmpty() {
					return nil, fmt.Errorf("Html parsing error: superfluous closing tag at %s", scanner.Location())
				}

				htmlEndPos := scanner.Position()

				closingTag, err := parseClosingTag(scanner)
				if err != nil {
					return nil, err
				}

				expected, ok := closeStack.Pop()
				if !ok || expected != closingTag {
					return nil, fmt.Errorf("Html parsing error: invalid closing tag '%s', expected '%s' at %s", closingTag, expected, scanner.Location())
				}

				if current != nil && current.htmlStart > 0 && current.htmlStart < htmlEndPos {
					current.InnerHtml = strings.TrimSpace(scanner.Slice(current.htmlStart, htmlEndPos))
				}

				current = current.Parent
				continue
			}

			tag, err := parsingOpenTag(parser, scanner)
			if err != nil {
				return nil, err
			}

			if current == nil {
				roots = append(roots, tag)
				tag.Parent = nil
			} else {
				tag.Parent = current
				current.Children = append(current.Children, tag)
			}

			if !tag.IsSelfClosing {
				current = tag
				closeStack.Push(tag.Name)
			}

		} else {
			contentStart := scanner.Position()
			scanner.ConsumeUntil(func(r rune) bool { return r == '<' })
			content := strings.TrimSpace(scanner.SliceFrom(contentStart))

			if content != "" && current != nil {
				current.InnerContent += content
			}
		}
	}

	if !closeStack.IsEmpty() {
		unclosed, _ := closeStack.Pop()
		return nil, fmt.Errorf("Html parsing error: unclosed tag '%s'", unclosed)
	}

	return roots, nil
}

func parseClosingTag(scanner *tools.Scanner) (string, error) {
	if !scanner.Match('<') {
		return "", fmt.Errorf("expected '<' at %s", scanner.Location())
	}

	if !scanner.Match('/') {
		return "", fmt.Errorf("expected '/' at %s", scanner.Location())
	}

	scanner.SkipWhitespace()

	tagName := scanner.ConsumeWhile(func(r rune) bool {
		return r != '>' && r != ' ' && r != '\t' && r != '\n' && r != '\r'
	})

	if tagName == "" {
		return "", fmt.Errorf("empty tag name at %s", scanner.Location())
	}

	scanner.SkipWhitespace()

	if !scanner.Match('>') {
		return "", fmt.Errorf("expected '>' at %s", scanner.Location())
	}

	return tagName, nil
}

func parsingOpenTag(parser *HtmlParser, scanner *tools.Scanner) (*HtmlTag, error) {
	startLine, startColumn := scanner.Line(), scanner.Column()
	if !scanner.Match('<') {
		return nil, fmt.Errorf("expected '<' at %s", scanner.Location())
	}

	scanner.SkipWhitespace()

	tagName := scanner.ConsumeWhile(func(r rune) bool {
		return r != '>' && r != '/' && r != ' ' && r != '\t' && r != '\n' && r != '\r'
	})

	if tagName == "" {
		return nil, fmt.Errorf("empty tag name at %s", scanner.Location())
	}

	handler, isExist := parser.customHandlers[tagName]

	if isExist {
		scanner.SetLocation(startLine, startColumn)
		return handler(scanner)
	}

	tag := &HtmlTag{
		Name:       tagName,
		Attributes: make([]HtmlAttribute, 0),
		Children:   make([]*HtmlTag, 0),
		Pos:        Position{Line: scanner.Line(), Column: scanner.Column()},
	}

	scanner.SkipWhitespace()

	for !scanner.EOF() && scanner.Ch != '>' && scanner.Ch != '/' {
		attr, err := parseAttribute(scanner)
		if err != nil {
			return nil, err
		}
		tag.Attributes = append(tag.Attributes, attr)
		scanner.SkipWhitespace()
	}

	if scanner.Ch == '/' {
		scanner.Take()
		tag.IsSelfClosing = true
	}

	if !scanner.Match('>') {
		return nil, fmt.Errorf("expected '>' at %s", scanner.Location())
	}

	tag.htmlStart = scanner.Position()

	return tag, nil
}

func parseAttribute(scanner *tools.Scanner) (HtmlAttribute, error) {
	scanner.SkipWhitespace()

	attrName := scanner.ConsumeWhile(func(r rune) bool {
		return r != '=' && r != '>' && r != '/' && r != ' ' && r != '\t' && r != '\n' && r != '\r'
	})

	if attrName == "" {
		return HtmlAttribute{}, fmt.Errorf("empty attribute name at %s", scanner.Location())
	}

	attr := HtmlAttribute{
		Name:         attrName,
		IsValueExist: false,
	}

	scanner.SkipWhitespace()

	if scanner.Ch == '=' {
		scanner.Take()
		scanner.SkipWhitespace()

		attr.IsValueExist = true

		if scanner.Ch == '"' || scanner.Ch == '\'' {
			quote := scanner.Ch
			scanner.Take()

			valueStart := scanner.Position()
			scanner.ConsumeUntil(func(r rune) bool { return r == quote })
			attr.Value = scanner.SliceFrom(valueStart)

			if !scanner.Match(quote) {
				return HtmlAttribute{}, fmt.Errorf("unclosed attribute value at %s", scanner.Location())
			}
		} else {
			attr.Value = scanner.ConsumeWhile(func(r rune) bool {
				return r != ' ' && r != '\t' && r != '\n' && r != '\r' && r != '>' && r != '/'
			})
		}
	}

	return attr, nil
}

func PrintHtmlTree(tag *HtmlTag) {
	printHtmlTreeRecursive(tag, 0)
}

func printHtmlTreeRecursive(tag *HtmlTag, depth int) {
	if tag == nil {
		return
	}

	indent := strings.Repeat("  ", depth)

	fmt.Printf("%s<%s", indent, tag.Name)

	for _, attr := range tag.Attributes {
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
