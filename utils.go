package main

import (
	"encoding/base64"
	"github.com/mdp/qrterminal/v3"
	qrcode "github.com/skip2/go-qrcode"
	"os"
	"reflect"
	"fmt"
	"crypto/md5"
	hashids "github.com/speps/go-hashids/v2"
	"strings"
)

func qrToURI(val string) string {
	png, _ := qrcode.Encode(val, qrcode.Medium, 256)
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
}

func qrToTerminal(val string) {
	qrterminal.GenerateHalfBlock(val, qrterminal.L, os.Stdout)
}

func ToMap(s interface{}, tagName string) map[string]interface{} {
	m := make(map[string]interface{})

	if s == nil {
		return m
	}

	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil
	}

	for i, t := 0, v.Type(); i < v.NumField(); i++ {
		fld := t.Field(i)
		tag := fld.Tag.Get(tagName)

		if tag != "" && tag != "-" {
			if fld.Type.Kind() == reflect.Struct {
				m[tag] = ToMap(v.Field(i).Interface(), tagName)
			} else {
				m[tag] = v.Field(i).Interface()
			}
		}
	}

	return m
}

func FromMap(s interface{}, m map[string]interface{}, tagName string) error {
	if s == nil {
		return fmt.Errorf("Target struct is nil: %s", s)
	}

	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	} else {
		return fmt.Errorf("Target must be a pointer to a struct")
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fld := v.Field(i)
		tag := t.Field(i).Tag.Get(tagName)

		if mv, ok := m[tag]; tag != "" && tag != "-" && ok {
			val := reflect.ValueOf(mv)

			if !fld.CanSet() {
				return fmt.Errorf("Cannot set %s field value", tag)
			}

			// TODO: Add support for nested structs
			// if fld.Type.Kind() == reflect.Struct {
			// 	FromMap(fld.Interface(), )
			// }

			if fld.Type() != val.Type() {
        		return fmt.Errorf("Field and value type don't match: %t", mv)
			}

			fld.Set(val)
		}
	}

	return nil
}

func HashString(s string) [md5.Size]byte {
	return md5.Sum([]byte(s))
}

func HashIDs(ids []int) string {
	hd := hashids.NewData()
	hd.Salt = "this is my salt"
	hd.MinLength = 28
	h, _ := hashids.NewWithData(hd)
	e, _ := h.Encode(ids)
	// fmt.Println(e)
	return e
}

func UnhashIDs(e string) []int {
	hd := hashids.NewData()
	hd.Salt = "this is my salt"
	hd.MinLength = 28
	h, _ := hashids.NewWithData(hd)
	d, _ := h.DecodeWithError(e)
	// fmt.Println(d)
	return d
}

func StripWhatsAppAt(s string) string {
	return strings.Split(s, "@")[0]
}
