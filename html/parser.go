package html

import (
	"fmt"
	"strings"

	"github.com/DilemaFixer/HtmlPuzzles/utils"
)

func ParseHtml(html string) ([]*HtmlTag, error) {
	scanner := utils.NewScanner(html)
	closeStack := utils.NewStack[string]()
	var roots []*HtmlTag
	var current *HtmlTag = nil

	var pendingMeta map[string]string

	for !scanner.EOF() {
		scanner.SkipWhitespace()
		if scanner.EOF() {
			break
		}

		if scanner.Current() == '<' {
			if scanner.Position()+4 <= scanner.Len() && scanner.Slice(scanner.Position(), scanner.Position()+4) == "<!--" {
				start := scanner.Position() + 4
				scanner.ConsumeUntilString("-->")
				if !scanner.MatchString("-->") {
					return nil, fmt.Errorf("Html parsing error: unclosed comment at %s", scanner.Location())
				}
				comment := strings.TrimSpace(scanner.Slice(start, scanner.Position()-3))

				// если это meta-коммент
				if strings.HasPrefix(comment, "meta:") {
					if pendingMeta == nil {
						pendingMeta = make(map[string]string)
					}
					metaStr := strings.TrimPrefix(comment, "meta:")
					parts := strings.SplitN(metaStr, "=", 2)
					if len(parts) == 2 {
						pendingMeta[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
					} else {
						pendingMeta[strings.TrimSpace(metaStr)] = ""
					}
				}

				continue
			}

			if scanner.Position()+9 <= scanner.Len() && scanner.Slice(scanner.Position(), scanner.Position()+9) == "<!DOCTYPE" {
				scanner.ConsumeUntil(func(r rune) bool { return r == '>' })
				if !scanner.Match('>') {
					return nil, fmt.Errorf("Html parsing error: unclosed DOCTYPE at %s", scanner.Location())
				}
				continue
			}

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

			tag, err := parsingOpenTag(scanner)
			if err != nil {
				return nil, err
			}

			if pendingMeta != nil {
				tag.Meta = pendingMeta
				pendingMeta = nil
			}

			if current == nil {
				roots = append(roots, tag)
			} else {
				tag.Parent = current
				current.Children = append(current.Children, tag)
			}

			if !tag.IsSelfClosing {
				current = tag
				closeStack.Push(tag.Name)
			}

		} else {
			// текстовый контент
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

func parseClosingTag(scanner *utils.Scanner) (string, error) {
	if !(scanner.Match('<') && scanner.Match('/')) {
		return "", fmt.Errorf("expected '</' at %s", scanner.Location())
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

func parsingOpenTag(scanner *utils.Scanner) (*HtmlTag, error) {
	startLine, startColumn := scanner.Line(), scanner.Column()
	if !scanner.Match('<') {
		return nil, fmt.Errorf("expected '<' at %s", scanner.Location())
	}

	scanner.SkipWhitespace()
	tagName := scanner.ConsumeWhile(func(r rune) bool {
		return r != '>' && r != '/' && !isWhitespace(r)
	})
	if tagName == "" {
		return nil, fmt.Errorf("empty tag name at %s", scanner.Location())
	}

	if tagName == "style" || tagName == "script" {
		for scanner.Current() != '>' && !scanner.EOF() {
			scanner.Take()
		}
		if !scanner.Match('>') {
			return nil, fmt.Errorf("expected '>' at %s", scanner.Location())
		}

		contentStart := scanner.Position()
		endTag := fmt.Sprintf("</%s>", tagName)
		scanner.ConsumeUntilString(endTag)
		content := scanner.Slice(contentStart, scanner.Position())
		scanner.MatchString(endTag)

		return &HtmlTag{
			Name:          tagName,
			Attributes:    make(map[string]HtmlAttribute),
			Children:      nil,
			InnerHtml:     content,
			IsSelfClosing: true,
			Pos:           Position{Line: startLine, Column: startColumn},
		}, nil
	}

	tag := &HtmlTag{
		Name:       tagName,
		Attributes: make(map[string]HtmlAttribute),
		Children:   make([]*HtmlTag, 0),
		Pos:        Position{Line: startLine, Column: startColumn},
	}

	scanner.SkipWhitespace()
	for !scanner.EOF() && scanner.Current() != '>' && scanner.Current() != '/' {
		attr, err := parseAttribute(scanner)
		if err != nil {
			return nil, err
		}
		tag.Attributes[attr.Name] = attr
		scanner.SkipWhitespace()
	}

	if scanner.Current() == '/' {
		scanner.Take()
		tag.IsSelfClosing = true
	}

	if !scanner.Match('>') {
		return nil, fmt.Errorf("expected '>' at %s", scanner.Location())
	}

	tag.htmlStart = scanner.Position()
	return tag, nil
}

func parseAttribute(scanner *utils.Scanner) (HtmlAttribute, error) {
	scanner.SkipWhitespace()
	attrName := scanner.ConsumeWhile(func(r rune) bool {
		return r != '=' && r != '>' && r != '/' && !isWhitespace(r)
	})

	if attrName == "" {
		return HtmlAttribute{}, fmt.Errorf("empty attribute name at %s", scanner.Location())
	}

	attr := HtmlAttribute{Name: attrName, IsValueExist: false}
	scanner.SkipWhitespace()

	if scanner.Current() == '=' {
		scanner.Take()
		scanner.SkipWhitespace()
		attr.IsValueExist = true

		if scanner.Current() == '"' || scanner.Current() == '\'' {
			quote := scanner.Current()
			scanner.Take()
			valueStart := scanner.Position()
			scanner.ConsumeUntil(func(r rune) bool { return r == quote })
			attr.Value = strings.TrimSpace(scanner.SliceFrom(valueStart))
			if !scanner.Match(quote) {
				return HtmlAttribute{}, fmt.Errorf("unclosed attribute value at %s", scanner.Location())
			}
		} else {
			attr.Value = strings.TrimSpace(scanner.ConsumeWhile(func(r rune) bool {
				return !isWhitespace(r) && r != '>' && r != '/'
			}))
		}
	}

	return attr, nil
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n' || r == '\r'
}
