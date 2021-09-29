package main

import (
	"encoding/base64"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/mdp/qrterminal/v3"
	"os"
	"reflect"
	// "fmt"
)

func qrToURI(val string) string {
	png, _ := qrcode.Encode(val, qrcode.Medium, 256)
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
}

func qrToTerminal(val string) {
	qrterminal.GenerateHalfBlock(val, qrterminal.L, os.Stdout)
}

func ToMap(s interface{}, tagName string) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	if s == nil {
		return m, nil
	}

	v := reflect.ValueOf(s)

    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }

    if v.Kind() != reflect.Struct {
        return nil, fmt.Errorf("ToMap only accepts structs; got %T", v)
    }

    for i, t := 0, v.Type(); i < v.NumField(); i++ {
    	fld := t.Field(i)
    	tag := fld.Tag.Get(tagName)

    	if tag != "" && tag != "-" {
            m[tag] = v.Field(i).Interface()
            if fld.Type.Kind() == reflect.Struct {
            	m[tag], _ = ToMap(v.Field(i).Interface(), tagName)
            } else {
            	m[tag] = v.Field(i).Interface()
            }
        }
    }

    return m, nil
}
