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
// Please reset output before Map()
func Map(lv lua.LValue, output interface{}) error {
	return NewMapper().Map(lv, output)
}

// Map maps the lua value to the given go pointer.
// Please reset output before Map()
func (m *Mapper) Map(lv lua.LValue, output interface{}) error {
	rv := reflect.ValueOf(output)
	if rv.Kind() != reflect.Ptr {
		return errors.New("output must be a pointer")
	}
	return m.MapValue(lv, rv.Elem())
}

// MapValue maps the lua value to go value.
// Please reset go value before MapValue()
func (m *Mapper) MapValue(lv lua.LValue, rv reflect.Value) error {
	if lv == lua.LNil {
		return nil // keep the old value
	}
	return m.mapNonNilValue(lv, rv)
}

func (m *Mapper) mapNonNilValue(lv lua.LValue, rv reflect.Value) error {
	assert.True(lv != lua.LNil) // lv is not *lua.LNilType
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
		return m.mapUint(lv, rv)
	case reflect.Uint8:
		return m.mapUint8(lv, rv)
	case reflect.Uint16:
		return m.mapUint16(lv, rv)
	case reflect.Uint32:
		return m.mapUint32(lv, rv)
	case reflect.Uint64:
		return m.mapUint64(lv, rv)
	case reflect.Uintptr:
		return TBI
	case reflect.Float32:
		return m.mapFloat32(lv, rv)
	case reflect.Float64:
		return m.mapFloat64(lv, rv)
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
		return m.mapInterface(lv, rv)
	case reflect.Map:
		return m.mapMap(lv, rv)
	case reflect.Ptr:
		return m.mapPtr(lv, rv)
	case reflect.Slice:
		return m.mapSlice(lv, rv)
	case reflect.String:
		return m.mapString(lv, rv)
	case reflect.Struct:
		return m.mapStruct(lv, rv)
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
	case *lua.LUserData:
		if b, ok := v.Value.(bool); ok {
			rv.SetBool(b)
			return nil
		}
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapInt(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Int)
	switch v := lv.(type) {
	case lua.LNumber:
		rv.SetInt(int64(v))
		return nil
	case *lua.LUserData:
		if n, ok := v.Value.(int); ok {
			rv.SetInt(int64(n))
			return nil
		}
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapInt8(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Int8)
	switch v := lv.(type) {
	case lua.LNumber:
		rv.SetInt(int64(v))
		return nil
	case *lua.LUserData:
		if n, ok := v.Value.(int8); ok {
			rv.SetInt(int64(n))
			return nil
		}
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapInt16(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Int16)
	switch v := lv.(type) {
	case lua.LNumber:
		rv.SetInt(int64(v))
		return nil
	case *lua.LUserData:
		if n, ok := v.Value.(int16); ok {
			rv.SetInt(int64(n))
			return nil
		}
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapInt32(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Int32)
	switch v := lv.(type) {
	case lua.LNumber:
		rv.SetInt(int64(v))
		return nil
	case *lua.LUserData:
		if n, ok := v.Value.(int32); ok {
			rv.SetInt(int64(n))
			return nil
		}
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapInt64(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Int64)
	switch v := lv.(type) {
	case lua.LNumber:
		rv.SetInt(int64(v))
		return nil
	case *lua.LUserData:
		if n, ok := v.Value.(int64); ok {
			rv.SetInt(int64(n))
			return nil
		}
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapUint(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Uint)
	switch v := lv.(type) {
	case lua.LNumber:
		rv.SetUint(uint64(v))
		return nil
	case *lua.LUserData:
		if n, ok := v.Value.(uint); ok {
			rv.SetUint(uint64(n))
			return nil
		}
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapUint8(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Uint8)
	switch v := lv.(type) {
	case lua.LNumber:
		rv.SetUint(uint64(v))
		return nil
	case *lua.LUserData:
		if n, ok := v.Value.(uint8); ok {
			rv.SetUint(uint64(n))
			return nil
		}
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapUint16(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Uint16)
	switch v := lv.(type) {
	case lua.LNumber:
		rv.SetUint(uint64(v))
		return nil
	case *lua.LUserData:
		if n, ok := v.Value.(uint16); ok {
			rv.SetUint(uint64(n))
			return nil
		}
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapUint32(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Uint32)
	switch v := lv.(type) {
	case lua.LNumber:
		rv.SetUint(uint64(v))
		return nil
	case *lua.LUserData:
		if n, ok := v.Value.(uint32); ok {
			rv.SetUint(uint64(n))
			return nil
		}
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapUint64(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Uint64)
	switch v := lv.(type) {
	case lua.LNumber:
		rv.SetUint(uint64(v))
		return nil
	case *lua.LUserData:
		if n, ok := v.Value.(uint64); ok {
			rv.SetUint(uint64(n))
			return nil
		}
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapFloat32(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Float32)
	switch v := lv.(type) {
	case lua.LNumber:
		rv.SetFloat(float64(v))
		return nil
	case *lua.LUserData:
		if f, ok := v.Value.(float32); ok {
			rv.SetFloat(float64(f))
			return nil
		}
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapFloat64(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Float64)
	switch v := lv.(type) {
	case lua.LNumber:
		rv.SetFloat(float64(v))
		return nil
	case *lua.LUserData:
		if f, ok := v.Value.(float64); ok {
			rv.SetFloat(f)
			return nil
		}
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapInterface(lv lua.LValue, rv reflect.Value) error {
	assert.True(lv != lua.LNil)
	assert.True(rv.Kind() == reflect.Interface)
	switch v := lv.(type) {
	case lua.LBool:
		rv.Set(reflect.ValueOf(bool(v)))
	case lua.LNumber:
		rv.Set(reflect.ValueOf(float64(v)))
	case lua.LString:
		rv.Set(reflect.ValueOf(string(v)))
	case *lua.LFunction:
		// ignore
	case *lua.LUserData:
		return m.mapLuaUserDataToGoValue(v, rv)
	// case *lua.LTThread: no such type
	case *lua.LTable:
		return m.mapLuaTableToGoInterface(v, rv)
	case *lua.LChannel:
		// ignore
	}
	return nil
}

func (m *Mapper) mapMap(lv lua.LValue, rv reflect.Value) error {
	assert.True(lv != lua.LNil)
	assert.True(rv.Kind() == reflect.Map)
	switch v := lv.(type) {
	case *lua.LTable:
		return m.mapLuaTableToGoMap(v, rv)
	case *lua.LUserData:
		return m.mapLuaUserDataToGoValue(v, rv)
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapPtr(lv lua.LValue, rv reflect.Value) error {
	assert.True(lv != lua.LNil)
	assert.True(rv.Kind() == reflect.Ptr)
	if ud, ok := lv.(*lua.LUserData); ok {
		return m.mapLuaUserDataToGoValue(ud, rv)
	}
	rv.Set(reflect.New(rv.Type().Elem()))
	return m.mapNonNilValue(lv, rv.Elem())
}

func (m *Mapper) mapSlice(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Slice)
	switch v := lv.(type) {
	case *lua.LTable:
		// Make a new slice and copy each element.
		tblLen := v.Len()
		rv.Set(reflect.MakeSlice(rv.Type(), tblLen, tblLen))
		return m.mapLuaTableToGoSlice(v, rv)
	case *lua.LUserData:
		return m.mapLuaUserDataToGoValue(v, rv)
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapLuaTableToGoSlice(tbl *lua.LTable, rv reflect.Value) error {
	assert.True(tbl != nil)
	assert.True(rv.Kind() == reflect.Slice)
	tblLen := tbl.Len()
	assert.True(rv.Len() >= tblLen)
	for i := 0; i < tblLen; i++ {
		if err := m.MapValue(tbl.RawGetInt(i+1), rv.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

func (m *Mapper) mapLuaUserDataToGoValue(ud *lua.LUserData, rv reflect.Value) error {
	udValue := ud.Value
	if udValue == nil {
		if canBeNil(rv) {
			rv.Set(reflect.Zero(rv.Type()))
			return nil
		}
		return &typeError{
			goType:                rv.Type(),
			luaType:               lua.LTUserData,
			isLuaUserDataValueNil: true,
		}
	}

	// must be the same type
	udValType := reflect.TypeOf(udValue)
	if udValType != rv.Type() {
		return &typeError{
			goType:               rv.Type(),
			luaType:              lua.LTUserData,
			luaUserDataValueType: udValType,
		}
	}
	rv.Set(reflect.ValueOf(udValue))
	return nil
}

func (m *Mapper) mapString(lv lua.LValue, rv reflect.Value) error {
	assert.True(lv != lua.LNil)
	assert.True(rv.Kind() == reflect.String)
	switch v := lv.(type) {
	case lua.LString:
		rv.SetString(string(v))
		return nil
	case *lua.LUserData:
		if s, ok := v.Value.(string); ok {
			rv.SetString(s)
			return nil
		}
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapStruct(lv lua.LValue, rv reflect.Value) error {
	assert.True(lv != lua.LNil)
	assert.True(rv.Kind() == reflect.Struct)
	switch v := lv.(type) {
	case *lua.LTable:
		return m.mapLuaTableToGoStruct(v, rv)
	case *lua.LUserData:
		return m.mapLuaUserDataToGoValue(v, rv)
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapLuaTableToGoStruct(tbl *lua.LTable, rv reflect.Value) error {
	assert.True(tbl != nil)
	assert.True(rv.Kind() == reflect.Struct)
	for i := 0; i < rv.NumField(); i++ {
		fldVal := rv.Field(i)
		if !fldVal.CanSet() {
			continue // unexported field
		}
		name := rv.Type().Field(i).Name // TODO: use tag
		if err := m.MapValue(tbl.RawGet(lua.LString(name)), fldVal); err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
	}
	return nil
}

func (m *Mapper) mapLuaTableToGoInterface(tbl *lua.LTable, rv reflect.Value) error {
	assert.True(tbl != nil)
	assert.True(rv.Kind() == reflect.Interface)
	maxn := tbl.MaxN()
	if maxn == 0 { // table -> map[interface{}]interface{}
		rv.Set(reflect.MakeMap(reflect.TypeOf(map[interface{}]interface{}{})))
		return m.mapLuaTableToGoMap(tbl, rv)
	} else { // array -> []interface{}
		rv.Set(reflect.MakeSlice(reflect.TypeOf([]interface{}{}), 0, maxn))
		return m.mapLuaTableToGoSlice(tbl, rv)
	}
}

func (m *Mapper) mapLuaTableToGoMap(tbl *lua.LTable, rv reflect.Value) error {
	assert.True(tbl != nil)
	assert.True(rv.Kind() == reflect.Map)
	mapType := rv.Type()
	keyType := mapType.Key()
	elemType := mapType.Elem()
	tbl.ForEach(func(lKey, lVal lua.LValue) {
		rvKeyPtr := reflect.New(keyType) // rvKeyPtr is a pointer to a new zero key
		rvKey := rvKeyPtr.Elem()
		if err := m.MapValue(lKey, rvKeyPtr.Elem()); err != nil { // TODO: MapValue() 应该只返回bool, 不要创建 error
			return // skip field if error
		}
		rvElemPtr := reflect.New(elemType)
		rvElem := rvElemPtr.Elem()
		if err := m.MapValue(lVal, rvElemPtr.Elem()); err != nil {
			return // skip field if error
		}
		rv.SetMapIndex(rvKey, rvElem)
	})
	return nil
}

// canBeNil reports whether its argument v can be nil.
// The nilable argument must be a chan, func, interface, map, pointer, or slice value.
func canBeNil(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Map,
		reflect.Ptr,
		reflect.UnsafePointer,
		reflect.Interface,
		reflect.Slice:
		return true
	}
	return false
}
