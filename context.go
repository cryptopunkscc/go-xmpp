package xmpp

import (
	"encoding/xml"
	"fmt"
	"io"
)

// Template defines a minimal interface for XML elements
type Template interface{}

// Context represents a space of known XML elements
type Context struct {
	templates       map[string]Template
	defaultTemplate Template
}

// NewContext instantiates a new element space with an optional default template
func NewContext(dt Template) *Context {
	return &Context{
		templates:       make(map[string]Template),
		defaultTemplate: dt,
	}
}

func (ctx *Context) Add(p Template) {
	key := ResolveName(p)
	if s := ResolveNamespace(p); s != "" {
		key = s + " " + key
	}

	ctx.templates[key] = p
}

func (ctx *Context) GetTemplate(start *xml.StartElement) Template {
	proto := ctx.lookupTemplate(start)

	if proto == nil {
		return nil
	}

	return clone(proto)
}

func (ctx *Context) SetDefaultTemplate(t Template) {
	ctx.defaultTemplate = t
}

// DecodeElement identifies and unmarshals an element
func (ctx *Context) DecodeElement(dec *xml.Decoder, start *xml.StartElement) (Template, error) {
	proto := ctx.GetTemplate(start)
	if proto == nil {
		return nil, fmt.Errorf("template for %s not found", start.Name.Local)
	}

	err := dec.DecodeElement(proto, start)
	if err != nil {
		return nil, err
	}

	return proto, nil
}

// Decode decodes the next element in the XML document
func (ctx *Context) Decode(dec *xml.Decoder) (Template, error) {
	for {
		// Fetch the next token
		token, err := dec.Token()
		if err != nil {
			return nil, err
		}

		switch typed := token.(type) {
		case xml.StartElement:
			return ctx.DecodeElement(dec, &typed)
		case xml.EndElement:
			return nil, ErrEndOfElement
		}
	}
}

// DecodeAll decodes all elements until the end of the current element
func (ctx *Context) DecodeAll(dec *xml.Decoder) ([]Template, error) {
	list := make([]Template, 0)

	for {
		item, err := ctx.Decode(dec)

		if err != nil {
			if (err == ErrEndOfElement) || (err == io.EOF) {
				return list, nil
			}
			return nil, err
		}

		list = append(list, item)
	}
}

func (ctx *Context) lookupTemplate(start *xml.StartElement) Template {
	key := start.Name.Local
	if start.Name.Space != "" {
		key = start.Name.Space + " " + key
	}

	// Check space+name match first
	if p, ok := ctx.templates[key]; ok {
		return p
	}

	// Check name-only match
	if p, ok := ctx.templates[start.Name.Local]; ok {
		return p
	}

	return ctx.defaultTemplate
}
