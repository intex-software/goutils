package inline

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrGivenNonStringKey     = errors.New("given object with non-string key")
	ErrMissingInlineField    = errors.New("target struct missing specified inline field")
	ErrNotGivenMutable       = errors.New("not given something which we can assign to")
	ErrNotStructHolder       = errors.New("holder is not a struct")
	ErrInlineNotRightMap     = errors.New("target's inline field is not a map[string]interface{}")
	ErrUnsettableInlineField = errors.New("target struct's inline field not assignable")
)

const (
	jsonTag   = "json"
	inlineTag = "inline"

	ErrMalformedJSON  = "given malformed JSON"
	ErrNotGivenStruct = "not given a struct in the raw stream"
)

// Contains reports whether v is present in s.
func Contains[S ~[]E, E comparable](s S, v E) bool {
	return Index(s, v) >= 0
}

// Index returns the index of the first occurrence of v in s,
// or -1 if not present.
func Index[S ~[]E, E comparable](s S, v E) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}

// UnmarshalJSON unmarshals a JSON object into a struct, with support for inline fields.
// Inline fields are fields that are not explicitly defined in the struct, but are instead
// stored in a map[string]interface{} field with the tag `json:"-,inline"`.
//
// Example:
//
//	type MyStruct struct {
//		ExplicitField string `json:"explicit"`
//		InlineField   map[string]interface{} `json:"-,inline"`
//	}
//
//	func (target *MyStruct) UnmarshalJSON(raw []byte) (err error) {
//		err = inline.UnmarshalJSON(target, raw)
//		if err == nil {
//			target.MetaData = nil
//		}
//		return
//	}
//
//	raw := []byte(`{"explicit":"value","extra":42}`)
//	var target MyStruct
//	err := json.UnmarshalJSON(&target, raw)
//	if err != nil {
//		panic(err)
//	}
//
//	fmt.Println(target.ExplicitField) // "value"
//	fmt.Println(target.InlineField)   // map[extra:42]
func InlineUnmarshalJSON(target any, raw []byte) (err error) {
	me := reflect.ValueOf(target)
	if me.Kind() != reflect.Ptr {
		return ErrNotGivenMutable
	}
	me = me.Elem()
	if me.Kind() != reflect.Struct {
		return ErrNotStructHolder
	}

	met := me.Type()
	fieldsLookup := make(map[string]int, met.NumField()-1)

	var inlineSink reflect.Value
	for i, length := 0, met.NumField(); i < length; i++ {
		sf := met.Field(i)
		if tag, ok := sf.Tag.Lookup(jsonTag); ok {
			sections := strings.Split(tag, ",")
			jsonName := sections[0]

			if jsonName != "-" {
				fieldsLookup[jsonName] = i
			} else if /* slices. */ Contains(sections[1:], inlineTag) {
				inlineSink = me.Field(i)
			}

		} else {
			fieldsLookup[sf.Name] = i
		}
	}
	if inlineSink.Kind() == reflect.Invalid {
		return ErrMissingInlineField
	}
	if inlineSink.Kind() != reflect.Map {
		return ErrInlineNotRightMap
	}
	if inlineSink.Type().Key().Kind() != reflect.String {
		return ErrInlineNotRightMap
	}

	inlineValueType := inlineSink.Type().Elem()
	if !inlineSink.CanSet() {
		return ErrUnsettableInlineField
	}

	dec := json.NewDecoder(bytes.NewReader(raw))
	if err := swallowRuneToken(dec, '{', ErrNotGivenStruct); err != nil {
		return err
	}
	for dec.More() {
		keyToken, err := dec.Token()
		if err != nil {
			return err
		}
		key, ok := keyToken.(string)
		if !ok {
			return ErrGivenNonStringKey
		}

		// dec.Token() skips over colons!

		var wantType reflect.Type
		if fieldIndex, ok := fieldsLookup[key]; ok {
			wantType = met.Field(fieldIndex).Type
		} else {
			wantType = inlineValueType
		}

		vvl := reflect.MakeSlice(reflect.SliceOf(wantType), 1, 1)
		vv := vvl.Index(0)
		err = dec.Decode(vv.Addr().Interface())
		if err != nil {
			return err
		}

		if fieldIndex, ok := fieldsLookup[key]; ok {
			me.Field(fieldIndex).Set(vv.Convert(met.Field(fieldIndex).Type))
		} else {
			kv := reflect.ValueOf(key)
			if inlineSink.IsNil() {
				inlineSink.Set(reflect.MakeMap(inlineSink.Type()))
			}
			inlineSink.SetMapIndex(kv, vv.Convert(inlineValueType))
		}
	}

	return swallowRuneToken(dec, '}', ErrMalformedJSON)
}

func swallowRuneToken(decoder *json.Decoder, expect rune, failExpect string) (err error) {
	t, err := decoder.Token()
	if err != nil {
		return
	}
	if t != json.Delim(expect) {
		err = fmt.Errorf("expected %q got %q: %s", expect, t, failExpect)
		return
	}
	return nil
}
