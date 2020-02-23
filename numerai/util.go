package numerai

import "reflect"

func structFields(s interface{}) []string {
	v := reflect.TypeOf(s)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	r := []string{}
	for i := 0; i < v.NumField(); i++ {
		r = append(r, v.Field(i).Name)
	}
	return r
}
