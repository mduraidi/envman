package internal

import (
	"reflect"
)

// StructToMap converts a struct to a map[string]string (exported fields only)
func StructToMap(s interface{}) map[string]string {
	result := make(map[string]string)
	v := reflect.ValueOf(s).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" { // unexported
			continue
		}
		val := v.Field(i)
		if val.Kind() == reflect.String {
			result[field.Name] = val.String()
		}
	}
	return result
}
