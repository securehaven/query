package query

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ParseFunc func(v string) (any, error)

type ValueFunc func(any) any

type field struct {
	name      string
	parseFunc ParseFunc
	valueFunc ValueFunc
}

func newField(name string, kind reflect.Kind) field {
	var parseFunc ParseFunc

	switch kind {
	case reflect.Bool:
		parseFunc = ParseBool
	case reflect.Int:
		parseFunc = ParseInt(0, 0)
	case reflect.Int8:
		parseFunc = ParseInt(0, 8)
	case reflect.Int16:
		parseFunc = ParseInt(0, 16)
	case reflect.Int32:
		parseFunc = ParseInt(0, 32)
	case reflect.Int64:
		parseFunc = ParseInt(0, 64)
	case reflect.Uint:
		parseFunc = ParseUint(0, 0)
	case reflect.Uint8:
		parseFunc = ParseUint(0, 8)
	case reflect.Uint16:
		parseFunc = ParseUint(0, 16)
	case reflect.Uint32:
		parseFunc = ParseUint(0, 32)
	case reflect.Uint64:
		parseFunc = ParseUint(0, 64)
	case reflect.Float32:
		parseFunc = ParseFloat(32)
	case reflect.Float64:
		parseFunc = ParseFloat(64)
	default:
		parseFunc = ParseString
	}

	return field{
		name:      name,
		parseFunc: parseFunc,
	}
}

var (
	ErrNoStruct = errors.New("expected struct as generic type")
)

func getFieldsFromStruct[S any]() ([]field, error) {
	var s S

	structType := reflect.TypeOf(s)
	kind := structType.Kind()

	if kind != reflect.Struct {
		return nil, fmt.Errorf("type %q is not a struct: %w", kind, ErrNoStruct)
	}

	return getFieldsFromReflectStruct(structType)
}

func getFieldsFromReflectStruct(st reflect.Type) ([]field, error) {
	fields := make([]field, 0, st.NumField())

	for i := 0; i < st.NumField(); i++ {
		structField := st.Field(i)
		name := getFieldNameFromStructField(structField)

		if len(name) == 0 {
			continue
		}

		structFieldType := structField.Type

		if structFieldType.Kind() == reflect.Pointer {
			structFieldType = structFieldType.Elem()
		}

		if structFieldType.Kind() != reflect.Struct {
			fields = append(fields, newField(name, structFieldType.Kind()))
			continue
		}

		childFields, err := getFieldsFromReflectStruct(structFieldType)

		if err != nil {
			return nil, err
		}

		fields = append(fields, newField(name, structFieldType.Kind()))

		for _, child := range childFields {
			fields = append(fields, field{
				name:      name + SeparatorSelector + child.name,
				parseFunc: child.parseFunc,
				valueFunc: child.valueFunc,
			})
		}
	}

	return fields, nil
}

func getFieldNameFromStructField(sf reflect.StructField) string {
	jsonTag := sf.Tag.Get("json")
	jsonTagTrimmed := strings.TrimSpace(jsonTag)

	if len(jsonTagTrimmed) == 0 {
		return sf.Name
	}

	jsonTagParts := strings.SplitN(jsonTagTrimmed, ",", 2)

	if len(jsonTagParts) == 0 {
		return sf.Name
	}

	jsonName := jsonTagParts[0]

	if len(jsonName) == 0 {
		return sf.Name
	}

	switch jsonName {
	case "-":
		return ""
	case "omitempty":
		return sf.Name
	}

	return jsonName
}
