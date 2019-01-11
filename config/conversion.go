package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	MapTag      = "map"
	TagRequired = "required"
)

// GetTagAttributes return the list of attribute of this tag
func GetTagAttributes(f reflect.StructField, tag string) []string {
	attributes, ok := f.Tag.Lookup(tag)
	if !ok {
		return nil
	}
	return strings.Split(attributes, ",")
}

func IsTagAttributePresent(f reflect.StructField, tag string, attr string) bool {
	lists := GetTagAttributes(f, tag)
	if lists == nil {
		return false
	}
	for _, l := range lists {
		if l == attr {
			return true
		}
	}

	return false
}

// StructToMap ...
func StructToMap(t interface{}) (string, map[string]interface{}, error) {
	m := make(map[string]interface{})
	typeT := reflect.TypeOf(t)
	typeName := typeT.String()
	// Erreur si le type n'est pas un pointeur
	if typeT.Kind() != reflect.Ptr {
		return "", nil, errors.New("t interface{} must be * to type but is " + typeName)
	}
	typeT = typeT.Elem()
	values := reflect.ValueOf(t)
	values = values.Elem()

	for i := 0; i < typeT.NumField(); i++ {
		ft := typeT.Field(i)
		fv := values.Field(i)
		// Regarde dans la map pour voire si on n'a le field
		switch fv.Type().Kind() {
		case reflect.String:
			m[ft.Name] = fv.String()
			break
		case reflect.Bool:
			m[ft.Name] = fv.Bool()
			break
		}
	}
	return typeT.Name(), m, nil
}

// MapToStruct convert a map to a Struct searching for the Type name for the key
func MapToStruct(m map[string]interface{}, t interface{}) error {
	typeT := reflect.TypeOf(t)
	typeName := typeT.String()
	// Erreur si le type n'est pas un pointeur
	if typeT.Kind() != reflect.Ptr {
		return errors.New("t interface{} must be * to type but is " + typeName)
	}
	typeT = typeT.Elem()
	typeName = typeT.Name()
	if v, ok := m[typeName]; ok {
		values := reflect.ValueOf(t)
		values = values.Elem()
		dataMap := v.(map[string]interface{})

		for i := 0; i < typeT.NumField(); i++ {
			ft := typeT.Field(i)
			fv := values.Field(i)
			// Regarde dans la map pour voire si on n'a le field
			if data, ok := dataMap[ft.Name]; ok {
				switch fv.Type().Kind() {
				case reflect.String:
					fv.SetString(data.(string))
					break
				case reflect.Bool:
					fv.SetBool(data.(bool))
					break
				}
			} else {
				// Si required throw error sinon continue
				if IsTagAttributePresent(ft, MapTag, TagRequired) {
					return errors.New("Field is required and not present in map " + ft.Name)
				}
			}
		}
		return nil
	}
	return fmt.Errorf("Key for Type %v does not exists", typeName)
}
