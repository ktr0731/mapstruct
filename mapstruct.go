package mapstruct

import (
	"errors"
	"reflect"
)

type mapper struct {
	c1, c2    reflect.Value
	isPointer bool
}

func Map(v1, v2 interface{}) (interface{}, error) {
	m, err := checkPrecondition(v1, v2)
	if err != nil {
		return nil, err
	}
	return m.mapStruct()
}

func checkPrecondition(v1, v2 interface{}) (*mapper, error) {
	rv1, rv2 := reflect.ValueOf(v1), reflect.ValueOf(v2)

	// check whether rv1 or rv2 is the zero value of empty interface type
	if !rv1.IsValid() || !rv2.IsValid() {
		return nil, errors.New("cannot map zero values")
	}

	if rv1.Type() != rv2.Type() {
		return nil, errors.New("not same types")
	}

	if rv1.Kind() == reflect.Interface || rv2.Kind() == reflect.Interface {
		return nil, errors.New("interface type is not supported yet")
	}

	m := &mapper{}

	// check whether rv1, 2 are pointer type
	m.obtain(rv1, rv2)

	if m.c1.Kind() != reflect.Struct || m.c2.Kind() != reflect.Struct {
		return nil, errors.New("cannot map not struct type")
	}

	return m, nil
}

func (m *mapper) obtain(rv1, rv2 reflect.Value) {
	m.c1 = m.obtainConcrete(rv1)
	m.c2 = m.obtainConcrete(rv2)
}

func (m *mapper) obtainConcrete(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Ptr {
		m.isPointer = true
		return m.obtainConcrete(v.Elem())
	}
	return v
}

func (m *mapper) mapStruct() (interface{}, error) {
	rv1, rv2 := m.c1, m.c2

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
	if m.isPointer {
		return res.Addr().Interface(), nil
	}
	return res.Interface(), nil
}
