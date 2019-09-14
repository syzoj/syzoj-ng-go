package problem

import (
	"bytes"
	"html"

	"github.com/beevik/etree"
	"gopkg.in/russross/blackfriday.v2"
	"github.com/microcosm-cc/bluemonday"
)

func renderStmt(buf *bytes.Buffer, w Warner, tok etree.Token) {
	switch obj := tok.(type) {
	case *etree.Element:
		for _, ch := range obj.Child {
			renderStmtNode(buf, w, ch)
		}
	}
}

func renderStmtNode(buf *bytes.Buffer, w Warner, tok etree.Token) {
	switch obj := tok.(type) {
	case *etree.Element:
		switch obj.Tag {
		case "Title":
			buf.WriteString("<h1>")
			renderText(buf, obj)
			buf.WriteString("</h1>")
		case "Description":
			buf.WriteString("<h2>Description</h2>")
			renderMarkdown(buf, obj)
		case "InputFormat":
			buf.WriteString("<h2>Input format</h2>")
			renderMarkdown(buf, obj)
		case "OutputFormat":
			buf.WriteString("<h2>Output format</h2>")
			renderMarkdown(buf, obj)
		case "Example":
			buf.WriteString("<h2>Example</h2>")
			renderMarkdown(buf, obj)
		case "LimitAndHint":
			buf.WriteString("<h2>Limit and hint</h2>")
			renderMarkdown(buf, obj)
		case "Tag":
		default:
			w.Warningf("Unrecognized statement tag: %s", obj.Tag)
		}
	}
}

var xssPolicy = bluemonday.UGCPolicy()
// TODO: customize and use latex
func renderMarkdown(buf *bytes.Buffer, tok etree.Token) {
	switch obj := tok.(type) {
	case *etree.CharData:
		log.Infof("char data: %#v", obj)
		res := blackfriday.Run([]byte(obj.Data))
		bytes := xssPolicy.SanitizeBytes(res)
		buf.Write(bytes)
	case *etree.Element:
		for _, child := range obj.Child {
			renderMarkdown(buf, child)
		}
	}
}

func renderText(buf *bytes.Buffer, tok etree.Token) {
	switch obj := tok.(type) {
	case *etree.CharData:
		buf.WriteString(html.EscapeString(obj.Data))
	case *etree.Element:
		for _, child := range obj.Child {
			renderText(buf, child)
		}
	}
}
