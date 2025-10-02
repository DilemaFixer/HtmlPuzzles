package html

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

func SerializeToString(tags []*HtmlTag) string {
	var sb strings.Builder
	for _, tag := range tags {
		renderTag(&sb, tag)
	}
	return sb.String()
}

func SerializeInWriter(w io.Writer, tags []*HtmlTag) error {
	for _, tag := range tags {
		if err := renderTagToWriter(w, tag); err != nil {
			return err
		}
	}
	return nil
}

func renderTagToWriter(w io.Writer, tag *HtmlTag) error {
	if tag == nil {
		return nil
	}

	if _, err := io.WriteString(w, "<"); err != nil {
		return err
	}
	if _, err := io.WriteString(w, tag.Name); err != nil {
		return err
	}

	keys := make([]string, 0, len(tag.Attributes))
	for k := range tag.Attributes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		attr := tag.Attributes[k]
		if attr.IsValueExist {
			if _, err := fmt.Fprintf(w, " %s=\"%s\"", attr.Name, attr.Value); err != nil {
				return err
			}
		} else {
			if _, err := io.WriteString(w, " "+attr.Name); err != nil {
				return err
			}
		}
	}

	if tag.IsSelfClosing {
		_, err := io.WriteString(w, "/>")
		return err
	}

	if _, err := io.WriteString(w, ">"); err != nil {
		return err
	}

	if tag.InnerContent != "" {
		if _, err := io.WriteString(w, tag.InnerContent); err != nil {
			return err
		}
	}

	for _, child := range tag.Children {
		if err := renderTagToWriter(w, child); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(w, "</"); err != nil {
		return err
	}
	if _, err := io.WriteString(w, tag.Name); err != nil {
		return err
	}
	if _, err := io.WriteString(w, ">"); err != nil {
		return err
	}

	return nil
}

func renderTag(sb *strings.Builder, tag *HtmlTag) {
	if tag == nil {
		return
	}

	sb.WriteString("<")
	sb.WriteString(tag.Name)

	keys := make([]string, 0, len(tag.Attributes))
	for k := range tag.Attributes {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		attr := tag.Attributes[k]
		if attr.IsValueExist {
			sb.WriteString(fmt.Sprintf(" %s=\"%s\"", attr.Name, attr.Value))
		} else {
			sb.WriteString(" " + attr.Name)
		}
	}

	if tag.IsSelfClosing {
		sb.WriteString("/>")
		return
	}

	sb.WriteString(">")
	if tag.InnerContent != "" {
		sb.WriteString(tag.InnerContent)
	}
	for _, child := range tag.Children {
		renderTag(sb, child)
	}
	sb.WriteString("</")
	sb.WriteString(tag.Name)
	sb.WriteString(">")
}
