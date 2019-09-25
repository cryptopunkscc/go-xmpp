package xmpp

import (
	"encoding/xml"
	"reflect"
	"strings"
)

// Identify returns an XMPP message name and namespace (can be empty)
func Identify(s interface{}) (id xml.Name) {
	// Check for explicit methods first
	type idName interface{ Name() string }
	type idSpace interface{ Namespace() string }
	if typed, ok := s.(idSpace); ok {
		id.Space = typed.Namespace()
	}
	if typed, ok := s.(idName); ok {
		id.Local = typed.Name()
		return
	}

	// Extract info from type's XMLName field tag
	id = extractXMLName(s)
	if id.Local != "" {
		return
	}

	// Finally, just use type's name in lowercase
	id = extractTypeName(s)

	return
}

// Extract type name via reflection (dereference a pointer if necessary)
func extractTypeName(s interface{}) xml.Name {
	typ := reflect.TypeOf(s)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return xml.Name{Local: strings.ToLower(typ.Name())}
}

// Extract XML element name and namespace from a structure
func extractXMLName(s interface{}) (id xml.Name) {
	typ := reflect.TypeOf(s)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	// Get the field
	field, ok := typ.FieldByName("XMLName")
	if !ok {
		return
	}

	// Get the XML tag
	tag, ok := field.Tag.Lookup("xml")
	if !ok || (len(tag) == 0) {
		return
	}

	// Split the tag
	parts := strings.Split(strings.Split(tag, ",")[0], " ")
	id.Local = parts[len(parts)-1]
	if len(parts) == 2 {
		id.Space = parts[0]
	}

	return
}
