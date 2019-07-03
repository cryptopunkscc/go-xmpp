package xmpp

import (
	"encoding/xml"
	"reflect"
	"strings"
)

func EncodeAll(enc *xml.Encoder, list []Template) error {
	var err error

	for _, item := range list {
		err = enc.Encode(item)
		if err != nil {
			return err
		}
	}

	return nil
}

func ResolveName(s interface{}) string {
	type namer interface {
		Name() string
	}

	switch t := s.(type) {
	case namer:
		return t.Name()
	}

	typ := reflect.TypeOf(s)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if sf, ok := typ.FieldByName("XMLName"); ok {
		if tag, ok := sf.Tag.Lookup("xml"); ok && len(tag) > 0 {
			parts := strings.Split(strings.Split(tag, ",")[0], " ")

			return parts[len(parts)-1]
		}
	}

	return strings.ToLower(typ.Name())
}

func ResolveNamespace(s interface{}) string {
	type namespacer interface {
		Namespace() string
	}

	// If the structure implements a Namespace() method, use it...
	switch t := s.(type) {
	case namespacer:
		return t.Namespace()
	}

	typ := reflect.TypeOf(s)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem() // turn the pointer type to the actual type
	}

	// If there's an XML tag with namespace in the structure, use it...
	if sf, ok := typ.FieldByName("XMLName"); ok {
		if tag, ok := sf.Tag.Lookup("xml"); ok && len(tag) > 0 {
			parts := strings.Split(strings.Split(tag, ",")[0], " ")
			if len(parts) == 2 {
				return parts[0]
			}
		}
	}

	// ... otherwise, there's no namespace
	return ""
}

func clone(i interface{}) interface{} {
	switch reflect.TypeOf(i).Kind() {
	case reflect.Struct:
		return reflect.ValueOf(i).Interface()

	case reflect.Ptr:
		ov := reflect.ValueOf(i).Elem()
		copy := reflect.New(ov.Type())
		copy.Elem().Set(ov)
		return copy.Interface()
	}

	return nil
}
