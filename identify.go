package xmpp

import (
	"reflect"
	"strings"
)

// Identify returns an XMPP message name and namespace (can be empty)
func Identify(s interface{}) (name string, space string) {
	// Check for explicit methods first
	type idName interface{ Name() string }
	type idSpace interface{ Namespace() string }
	if typed, ok := s.(idName); ok {
		name = typed.Name()
	}
	if typed, ok := s.(idSpace); ok {
		space = typed.Namespace()
	}

	// Extract info from type's XMLName field tag
	xmlname, xmlspace := extractXMLName(s)
	if name == "" {
		name = xmlname
	}
	if space == "" {
		space = xmlspace
	}

	// Finally, just use type's name in lowercase
	if name == "" {
		name = extractTypeName(s)
	}

	return
}

// Extract type name via reflection (dereference a pointer if necessary)
func extractTypeName(s interface{}) string {
	typ := reflect.TypeOf(s)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	return strings.ToLower(typ.Name())
}

// Extract XML element name and namespace from a structure
func extractXMLName(s interface{}) (name string, space string) {
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
	name = parts[len(parts)-1]
	if len(parts) == 2 {
		space = parts[0]
	}

	return
}