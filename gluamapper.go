// gluamapper provides an easy way to map GopherLua tables to Go structs.
package gluamapper

import (
	"errors"
	"fmt"
	"reflect"

	assert "github.com/arl/assertgo"
	"github.com/yuin/gopher-lua"
)

// Mapper maps a lua table to a Go struct pointer.
type Mapper struct {
	// A struct tag name for lua table keys.
	TagName string
}

// NewMapper returns a new mapper.
func NewMapper() *Mapper {
	return &Mapper{}
}

// Map maps the lua value to the given go pointer with default options.
func Map(lv lua.LValue, output interface{}) error {
	return NewMapper().Map(lv, output)
}

// Map maps the lua value to the given go pointer.
func (m *Mapper) Map(lv lua.LValue, output interface{}) error {
	rv := reflect.ValueOf(output)
	if rv.Kind() != reflect.Ptr {
		return errors.New("output must be a pointer")
	}
	return m.MapValue(lv, rv.Elem())
}

func (m *Mapper) MapValue(lv lua.LValue, rv reflect.Value) error {
	if _, ok := lv.(*lua.LNilType); ok {
		return nil
	}

	TBI := errors.New("to be implemented")
	switch rv.Kind() {
	case reflect.Bool:
		return m.mapBool(lv, rv)
	case reflect.Int:
		return m.mapInt(lv, rv)
	case reflect.Int8:
		return m.mapInt8(lv, rv)
	case reflect.Int16:
		return m.mapInt16(lv, rv)
	case reflect.Int32:
		return m.mapInt32(lv, rv)
	case reflect.Int64:
		return m.mapInt64(lv, rv)
	case reflect.Uint:
		return TBI
	case reflect.Uint8:
		return TBI
	case reflect.Uint16:
		return TBI
	case reflect.Uint32:
		return TBI
	case reflect.Uint64:
		return TBI
	case reflect.Uintptr:
		return TBI
	case reflect.Float32:
		return TBI
	case reflect.Float64:
		return TBI
	case reflect.Complex64:
		return TBI
	case reflect.Complex128:
		return TBI
	case reflect.Array:
		return TBI
	case reflect.Chan:
		return TBI
	case reflect.Func:
		return TBI
	case reflect.Interface:
		return TBI
	case reflect.Map:
		return TBI
	case reflect.Ptr:
		return TBI
	case reflect.Slice:
		return TBI
	case reflect.String:
		return TBI
	case reflect.Struct:
		return TBI
	case reflect.UnsafePointer:
		return TBI
	}
	return fmt.Errorf("unsupported type: %s", rv.Kind())
}

func (m *Mapper) mapBool(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Bool)
	switch v := lv.(type) {
	case lua.LBool:
		rv.SetBool(bool(v))
		return nil
	case (*lua.LUserData):
		if b, ok := v.Value.(bool); ok {
			rv.SetBool(b)
			return nil
		}
	}
	return typeError("bool", lv)
}

func (m *Mapper) mapInt(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Int)
	switch v := lv.(type) {
	case lua.LNumber:
		rv.SetInt(int64(v))
		return nil
	case (*lua.LUserData):
		if n, ok := v.Value.(int); ok {
			rv.SetInt(int64(n))
			return nil
		}
	}
	return typeError("int", lv)
}

func (m *Mapper) mapInt8(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Int8)
	switch v := lv.(type) {
	case lua.LNumber:
		rv.SetInt(int64(v))
		return nil
	case (*lua.LUserData):
		if n, ok := v.Value.(int8); ok {
			rv.SetInt(int64(n))
			return nil
		}
	}
	return typeError("int8", lv)
}

func (m *Mapper) mapInt16(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Int16)
	switch v := lv.(type) {
	case lua.LNumber:
		rv.SetInt(int64(v))
		return nil
	case (*lua.LUserData):
		if n, ok := v.Value.(int16); ok {
			rv.SetInt(int64(n))
			return nil
		}
	}
	return typeError("int16", lv)
}

func (m *Mapper) mapInt32(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Int32)
	switch v := lv.(type) {
	case lua.LNumber:
		rv.SetInt(int64(v))
		return nil
	case (*lua.LUserData):
		if n, ok := v.Value.(int32); ok {
			rv.SetInt(int64(n))
			return nil
		}
	}
	return typeError("int32", lv)
}

func (m *Mapper) mapInt64(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Int64)
	switch v := lv.(type) {
	case lua.LNumber:
		rv.SetInt(int64(v))
		return nil
	case (*lua.LUserData):
		if n, ok := v.Value.(int64); ok {
			rv.SetInt(int64(n))
			return nil
		}
	}
	return typeError("int64", lv)
}

func typeError(expectedTypeName string, lv lua.LValue) error {
	if ud, ok := lv.(*lua.LUserData); ok {
		val := reflect.ValueOf(ud.Value)
		typ := reflect.Indirect(val).Type()
		return fmt.Errorf("%s expected but got lua user data of %s", expectedTypeName, typ)
	}
	return fmt.Errorf("%s expected but got lua %s", expectedTypeName, lv.Type())
}
