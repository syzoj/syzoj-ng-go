package problem

import (
	"bytes"
	"html"

	"github.com/microcosm-cc/bluemonday"
	"github.com/syzoj/syzoj-ng-go/lib/xml"
	"gopkg.in/russross/blackfriday.v2"
)

func renderStmt(buf *bytes.Buffer, w Warner, tok xml.Token) {
	switch obj := tok.(type) {
	case *xml.Element:
		for _, ch := range obj.Child {
			renderStmtNode(buf, w, ch)
		}
	}
}

func renderStmtNode(buf *bytes.Buffer, w Warner, tok xml.Token) {
	switch obj := tok.(type) {
	case *xml.Element:
		switch obj.Name.Local {
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
			w.Warningf("Unrecognized statement tag: %s", obj.Name.Local)
		}
	}
}

var xssPolicy = bluemonday.UGCPolicy()

// TODO: customize and use latex
func renderMarkdown(buf *bytes.Buffer, tok xml.Token) {
	res := blackfriday.Run(xml.GetCharDataBytes(tok))
	bytes := xssPolicy.SanitizeBytes(res)
	buf.Write(bytes)
}

func renderText(buf *bytes.Buffer, tok xml.Token) {
	buf.WriteString(html.EscapeString(string(xml.GetCharDataBytes(tok))))
}
