package egoscale

import (
	"fmt"
	"reflect"
)

// Copy copies the value from from into to. The type of to must be convertible into the type of from.
func Copy(to, from interface{}) error {
	tt := reflect.TypeOf(to).Elem()
	ft := reflect.TypeOf(from).Elem()
	if !ft.ConvertibleTo(tt) {
		return fmt.Errorf("cannot convert %s into %s", tt.Name(), ft.Name())
	}
	tv := reflect.ValueOf(to).Elem()
	fv := reflect.ValueOf(from).Elem()
	tv.Set(fv.Convert(tt))
	return nil
}
