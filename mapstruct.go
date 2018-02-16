package mapstruct

import (
	"errors"
	"reflect"
)

func Map(v1, v2 interface{}) (interface{}, error) {
	rv1, rv2, err := checkPrecondition(v1, v2)
	if err != nil {
		return nil, err
	}
	return mapStruct(rv1, rv2)
}

func checkPrecondition(v1, v2 interface{}) (reflect.Value, reflect.Value, error) {
	rv1, rv2 := reflect.ValueOf(v1), reflect.ValueOf(v2)

	// check whether rv1 or rv2 is the zero value of empty interface type
	if !rv1.IsValid() || !rv2.IsValid() {
		return rv1, rv2, errors.New("cannot map zero values")
	}

	if rv1.Type() != rv2.Type() {
		return rv1, rv2, errors.New("not same types")
	}

	if rv1.Kind() == reflect.Interface || rv2.Kind() == reflect.Interface {
		return rv1, rv2, errors.New("interface type is not supported yet")
	}

	// check whether rv1, 2 are pointer type

	if rv1.Kind() != reflect.Struct || rv2.Kind() != reflect.Struct {
		return rv1, rv2, errors.New("cannot map not sturuct type")
	}

	return rv1, rv2, nil
}

func mapStruct(rv1, rv2 reflect.Value) (interface{}, error) {
	resp := reflect.New(rv1.Type())
	res := reflect.Indirect(resp)

	// check each value and if the value is not empty, assign it to "to"'s field
	for i := 0; i < rv1.NumField(); i++ {
		to, from, resf := rv1.Field(i), rv2.Field(i), res.Field(i)

		// if from is the zero value, skip this field
		if !from.IsValid() {
			continue
		}

		switch from.Kind() {
		case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr:
			if from.IsNil() {
				continue
			}
		case reflect.Invalid:
			continue
		default:
			// primitive or struct

			pv := resf.Addr()

			// skip unexported fields
			if !pv.Elem().CanSet() {
				continue
			}

			// if the value of from is the zero value, ignore it
			if reflect.Zero(from.Type()).Interface() == from.Interface() {
				pv.Elem().Set(to)
			} else {
				pv.Elem().Set(from)
			}
		}
	}
	return res.Interface(), nil
}
