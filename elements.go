package xmpp

import (
	"encoding/xml"
	"reflect"
)

var elements map[xml.Name]interface{}

func AddElement(e interface{}) {
	id := Identify(e)
	if elements == nil {
		elements = make(map[xml.Name]interface{})
	}
	elements[id] = e
}

func GetElement(name xml.Name) interface{} {
	if e := lookupTemplate(name); e != nil {
		return cloneElement(e)
	}

	return &Generic{}
}

func lookupTemplate(name xml.Name) interface{} {
	if e, ok := elements[name]; ok {
		return e
	}
	if e, ok := elements[xml.Name{Local: name.Local}]; ok {
		return e
	}
	return nil
}

func cloneElement(i interface{}) interface{} {
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
