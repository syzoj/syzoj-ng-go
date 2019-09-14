// xml is a simple wrapper for encoding/xml that provides tree-oridented functions.
package xml

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
)

var ErrXML = errors.New("Invalid XML")

// Token can be one of the types: *Element, CharData, Comment, ProcInst, Directive.
type Token interface {
	EncodeTo(*xml.Encoder) error
}

// Some re-exports for convenience
type Encoder = xml.Encoder
type Decoder = xml.Decoder
type Name = xml.Name

// xml.CharData
type CharData []byte

func (c CharData) EncodeTo(e *xml.Encoder) error { return e.EncodeToken(xml.CharData(c)) }

// xml.Comment
type Comment []byte

func (c Comment) EncodeTo(e *xml.Encoder) error { return e.EncodeToken(xml.Comment(c)) }

// xml.ProcInst
type ProcInst struct {
	Target string
	Inst   []byte
}

func (p ProcInst) EncodeTo(e *xml.Encoder) error { return e.EncodeToken(xml.ProcInst(p)) }

// xml.Directive
type Directive []byte

func (d Directive) EncodeTo(e *xml.Encoder) error { return e.EncodeToken(xml.Directive(d)) }

// Like xml.Element but with children
type Element struct {
	Name  xml.Name
	Attr  []xml.Attr
	Child []Token
}

func (e *Element) EncodeTo(en *xml.Encoder) error {
	if err := en.EncodeToken(xml.StartElement{Name: e.Name, Attr: e.Attr}); err != nil {
		return err
	}
	for _, ch := range e.Child {
		if err := ch.EncodeTo(en); err != nil {
			return err
		}
	}
	if err := en.EncodeToken(xml.EndElement{Name: e.Name}); err != nil {
		return err
	}
	return nil
}

// Decodes any token except EndElement.
func DecodeDecoder(d *xml.Decoder) (Token, error) {
	tok, err := d.Token()
	if err != nil {
		return nil, err
	}
	return decodeToken(d, tok)
}

func decodeToken(d *xml.Decoder, tok xml.Token) (Token, error) {
	switch obj := tok.(type) {
	case xml.StartElement:
		return DecodeElementDecoder(d, &obj)
	case xml.EndElement:
		return nil, ErrXML
	case xml.CharData:
		return CharData(obj.Copy()), nil
	case xml.Comment:
		return Comment(obj.Copy()), nil
	case xml.ProcInst:
		return ProcInst(obj.Copy()), nil
	case xml.Directive:
		return Directive(obj.Copy()), nil
	default:
		return nil, ErrXML
	}
}

// Decodes an element. If st is nil, it looks for the first element and skips anything else.
func DecodeElementDecoder(d *xml.Decoder, st *xml.StartElement) (*Element, error) {
	if st == nil {
		for {
			tok, err := d.Token()
			if err != nil {
				return nil, err
			}
			if se, ok := tok.(xml.StartElement); ok {
				st = &se
				break
			} else if _, ok = tok.(*xml.EndElement); ok {
				return nil, ErrXML
			}
		}
	}

	e := &Element{Name: st.Name, Attr: st.Attr}
	for {
		tok, err := d.Token()
		if err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return e, err
		}
		if obj, ok := tok.(*xml.EndElement); ok {
			if obj.Name != e.Name {
				return e, ErrXML
			}
			break
		}
		tokd, err := decodeToken(d, tok)
		if err != nil {
			return e, nil
		}
		e.Child = append(e.Child, tokd)
	}
	return e, nil
}

// Decodes an element from a reader.
func DecodeFrom(r io.Reader) (*Element, error) {
	return DecodeElementDecoder(xml.NewDecoder(r), nil)
}

// Decodes an element from bytes.
func DecodeFromBytes(r []byte) (*Element, error) {
	return DecodeFrom(bytes.NewBuffer(r))
}

// Encodes a token into a writer.
func EncodeTo(tok Token, w io.Writer) error {
	enc := xml.NewEncoder(w)
	if err := tok.EncodeTo(enc); err != nil {
		return err
	}
	if err := enc.Flush(); err != nil {
		return err
	}
	return nil
}

// Encodes a token into bytes.
func EncodeToBytes(tok Token) []byte {
	buf := &bytes.Buffer{}
	EncodeTo(tok, buf)
	return buf.Bytes()
}

// Some helper functions

// Find a direct child with given tag.
func (e *Element) SelectElement(tagName string) *Element {
	for _, ch := range e.Child {
		if el, ok := ch.(*Element); ok {
			if el.Name.Local == tagName {
				return el
			}
		}
	}
	return nil
}

// Gets an attribute.
func (e *Element) SelectAttrDefault(name string, def string) string {
	for _, attr := range e.Attr {
		if attr.Name.Local == name {
			return attr.Value
		}
	}
	return def
}

func (e *Element) CreateAttr(name string, val string) {
	e.Attr = append(e.Attr, xml.Attr{
		Name:  xml.Name{Local: name},
		Value: val,
	})
}

// Iterate through the subtree.
func (e *Element) TraverseSubtree(f func(*Element)) {
	f(e)
	for _, ch := range e.Child {
		if el, ok := ch.(*Element); ok {
			el.TraverseSubtree(f)
		}
	}
}

// Extract char data from token.
func GetCharData(tok Token, w io.Writer) error {
	switch obj := tok.(type) {
	case CharData:
		_, err := w.Write(obj)
		return err
	case *Element:
		for _, ch := range obj.Child {
			if err := GetCharData(ch, w); err != nil {
				return err
			}
		}
		return nil
	default:
		return nil
	}
}

// Extract char data from token in bytes.
func GetCharDataBytes(tok Token) []byte {
	buf := &bytes.Buffer{}
	GetCharData(tok, buf)
	return buf.Bytes()
}
