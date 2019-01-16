package conv

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strconv"
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

// GetUnderlyingType retourne informations about the type
// behind the interface{} and the ptr ( must be )
func GetUnderlyingType(t interface{}) (string, reflect.Type, reflect.Value, error) {
	typeT := reflect.TypeOf(t)
	typeName := typeT.Name()
	values := reflect.ValueOf(t)
	if typeT.Kind() != reflect.Ptr {
		return typeName, typeT, values, &NotPointerError{
			Type: typeName,
		}
	}
	typeT = typeT.Elem()
	return typeT.Name(), typeT, values.Elem(), nil
}

func getFieldQuery(elem string, sf *reflect.StructField, fv *reflect.Value) error {
	switch fv.Type().Kind() {
	case reflect.String:
		fv.SetString(elem)
	case reflect.Int:
		i, err := strconv.Atoi(elem)
		if err != nil {
			return err
		}
		fv.SetInt(int64(i))
	}

	return nil
}

// QueryToStruct get the value of the struct from the value of the query
func QueryToStruct(q map[string][]string, t interface{}) error {
	_, typeT, values, err := GetUnderlyingType(t)
	if err != nil {
		return err
	}
	for i := 0; i < typeT.NumField(); i++ {
		ft := typeT.Field(i)
		fv := values.Field(i)
		if elems, ok := q[ft.Name]; ok {
			if fv.Type().Kind() == reflect.Array {

			} else {
				if len(elems) == 1 {
					err = getFieldQuery(elems[0], &ft, &fv)
					if err != nil {
						return err
					}
				} else {
					// erreur field unique
				}
			}
		} else {
			// valide si required
		}
	}
	return nil
}

func addField(query []string, sf *reflect.StructField, fv *reflect.Value) ([]string, error) {
	var s string
	switch fv.Type().Kind() {
	case reflect.String:
		s = fv.String()
	case reflect.Int:
		s = fmt.Sprintf("%v", fv.Int())
	default:
		return query, nil
	}
	query = append(query, sf.Name+"="+s)
	return query, nil
}

// StructToQuery add a struct to an http query with reflection
func StructToQuery(t interface{}) (string, error) {
	_, typeT, values, err := GetUnderlyingType(t)
	if err != nil {
		return "", err
	}
	var elems []string
	for i := 0; i < typeT.NumField(); i++ {
		ft := typeT.Field(i)
		fv := values.Field(i)
		if fv.Type().Kind() == reflect.Array {
		} else {
			elems, err = addField(elems, &ft, &fv)
			if err != nil {
				return "", err
			}
		}
	}
	// Generate the query string
	elemsStr := strings.Join(elems, "&")
	query := "?" + elemsStr

	return query, nil
}

// AddStructToMap add a new structure with is field to a map
// only one instance of each struct can be store , unless using
// an array or map
func AddStructToMap(m map[string]interface{}, t interface{}) error {
	name, data, err := StructToMap(t)
	if err != nil {
		return err
	}
	// Si la clé existe deja throw une erreur aussi
	if _, ok := m[name]; ok {
		return &KeyError{
			Name:   name,
			Status: AlreadySet,
		}
	}
	m[name] = data
	return nil
}

// StructToMap convert a struct to a map with the FieldName as key
// and there value in the interface{}
func StructToMap(t interface{}) (string, map[string]interface{}, error) {
	m := make(map[string]interface{})
	_, typeT, values, err := GetUnderlyingType(t)
	if err != nil {
		return "", nil, err
	}

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
		case reflect.Int:
			m[ft.Name] = int(fv.Int())
		}
	}
	return typeT.Name(), m, nil
}

// FindStructMap try to find if the map contains one entry for is type
// is so extract it to the interface
func FindStructMap(m map[string]interface{}, t interface{}) error {
	// Get le nom du type et regarde s'il est present dans ma map
	name, _, _, err := GetUnderlyingType(t)
	if err != nil {
		return err
	}
	if d, ok := m[name]; ok {
		// essaye de cast interface en Ctx
		if mdata, ok := d.(map[string]interface{}); ok {
			return MapToStruct(mdata, t)
		} else {
			typeT := reflect.TypeOf(mdata)
			return &BadTypeError{
				GotType:    typeT.Name(),
				WantedType: "map[string]interface{}",
			}
		}
	}
	return &KeyError{
		Name:   name,
		Status: NotFound,
	}
}

// MapToStruct convert a map to a Struct searching for the Type name for the key
func MapToStruct(m map[string]interface{}, t interface{}) error {
	_, typeT, values, err := GetUnderlyingType(t)
	if err != nil {
		return err
	}

	for i := 0; i < typeT.NumField(); i++ {
		ft := typeT.Field(i)
		fv := values.Field(i)
		// Regarde dans la map pour voire si on n'a le field
		if data, ok := m[ft.Name]; ok {
			switch fv.Type().Kind() {
			case reflect.String:
				fv.SetString(data.(string))
				break
			case reflect.Bool:
				fv.SetBool(data.(bool))
				break
			case reflect.Int:
				if di, ok := data.(int); ok {
					fv.SetInt(int64(di))
				} else if di, ok := data.(float64); ok {
					fv.SetInt(int64(di))
				}
				break
			default:

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

type SerializeExtension string

const (
	ExtJSON SerializeExtension = "json"
	ExtGLOB SerializeExtension = "data"
)

// Save save the map in the desire format determine with the
// file extensions name
func Save(filePath string, m map[string]interface{}) error {
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, b, 0644)
}

// Load the map in te desire format determine with the file extensions name
func Load(filePath string) (map[string]interface{}, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	d := make(map[string]interface{})
	return d, json.Unmarshal(b, &d)
}
