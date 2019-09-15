package xmpp

import (
	"encoding/xml"
	"reflect"
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
