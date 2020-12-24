package gluamapper

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	assert "github.com/arl/assertgo"
	"github.com/yuin/gopher-lua"
)

var (
	OutputValueIsNilError = errors.New("output value is nil")
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

// NewMapperWithTagName returns a new mapper with tag name.
func NewMapperWithTagName(tagName string) *Mapper {
	return &Mapper{
		TagName: tagName,
	}
}

// Map maps the lua value to the given go pointer.
func (m *Mapper) Map(lv lua.LValue, output interface{}) error {
	rv := reflect.ValueOf(output)
	if rv.Kind() != reflect.Ptr {
		return &OutputIsNotAPointerError{outputValue: rv}
	}
	return m.MapValue(lv, rv.Elem())
}

// MapValue maps the lua value to go value.
func (m *Mapper) MapValue(lv lua.LValue, rv reflect.Value) error {
	if lv != lua.LNil {
		return m.mapNonNilValue(lv, rv)
	}

	// do not call rv.Type() if rv is zero Value
	if rv.IsValid() {
		rv.Set(reflect.Zero(rv.Type()))
		return nil
	}
	return OutputValueIsNilError
}

func (m *Mapper) mapNonNilValue(lv lua.LValue, rv reflect.Value) error {
	assert.True(lv != lua.LNil) // lv is not *lua.LNilType
	TBI := errors.New("to be implemented")
	switch rv.Kind() {
	case reflect.Invalid:
		return OutputValueIsNilError
	case reflect.Bool:
		return mapBool(lv, rv)
	case reflect.Int:
		return mapInt(lv, rv)
	case reflect.Int8:
		return mapInt8(lv, rv)
	case reflect.Int16:
		return mapInt16(lv, rv)
	case reflect.Int32:
		return mapInt32(lv, rv)
	case reflect.Int64:
		return mapInt64(lv, rv)
	case reflect.Uint:
		return mapUint(lv, rv)
	case reflect.Uint8:
		return mapUint8(lv, rv)
	case reflect.Uint16:
		return mapUint16(lv, rv)
	case reflect.Uint32:
		return mapUint32(lv, rv)
	case reflect.Uint64:
		return mapUint64(lv, rv)
	case reflect.Uintptr:
		return TBI
	case reflect.Float32:
		return mapFloat32(lv, rv)
	case reflect.Float64:
		return mapFloat64(lv, rv)
	case reflect.Complex64:
		return TBI
	case reflect.Complex128:
		return TBI
	case reflect.Array: // TODO
		return TBI
	case reflect.Chan:
		return TBI
	case reflect.Func:
		return TBI
	case reflect.Interface:
		return mapInterface(lv, rv)
	case reflect.Map:
		return m.mapMap(lv, rv)
	case reflect.Ptr:
		return m.mapPtr(lv, rv)
	case reflect.Slice:
		return m.mapSlice(lv, rv)
	case reflect.String:
		return mapString(lv, rv)
	case reflect.Struct:
		return m.mapStruct(lv, rv)
	case reflect.UnsafePointer:
		return TBI
	}
	return fmt.Errorf("unsupported type: %s", rv.Kind())
}

func (m *Mapper) mapMap(lv lua.LValue, rv reflect.Value) error {
	assert.True(lv != lua.LNil)
	assert.True(rv.Kind() == reflect.Map)
	switch v := lv.(type) {
	case *lua.LTable:
		return m.mapLuaTableToGoMap(v, rv)
	case *lua.LUserData:
		return mapLuaUserDataToGoValue(v, rv)
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapPtr(lv lua.LValue, rv reflect.Value) error {
	assert.True(lv != lua.LNil)
	assert.True(rv.Kind() == reflect.Ptr)
	if ud, ok := lv.(*lua.LUserData); ok {
		return mapLuaUserDataToGoValue(ud, rv)
	}
	elemPtr := reflect.New(rv.Type().Elem())
	if err := m.mapNonNilValue(lv, elemPtr.Elem()); err != nil {
		return err
	}
	rv.Set(elemPtr)
	return nil
}

func (m *Mapper) mapSlice(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Slice)
	switch v := lv.(type) {
	case *lua.LTable:
		return m.mapLuaTableToGoSlice(v, rv)
	case *lua.LUserData:
		return mapLuaUserDataToGoValue(v, rv)
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapLuaTableToGoSlice(tbl *lua.LTable, rv reflect.Value) error {
	assert.True(tbl != nil)
	assert.True(rv.Kind() == reflect.Slice)
	tblLen := tbl.Len()
	rvCap := rv.Cap()
	if rvCap < tblLen {
		// reset to a new slice if need more capacity
		rv.Set(reflect.MakeSlice(rv.Type(), tblLen, tblLen))
	} else if rv.Len() != tblLen {
		// set len if capacity is large enough
		rv.SetLen(tblLen)
	}

	for i := 0; i < tblLen; i++ {
		if err := m.MapValue(tbl.RawGetInt(i+1), rv.Index(i)); err != nil {
			return err
		}
	}
	return nil
}

func (m *Mapper) mapStruct(lv lua.LValue, rv reflect.Value) error {
	assert.True(lv != lua.LNil)
	assert.True(rv.Kind() == reflect.Struct)
	switch v := lv.(type) {
	case *lua.LTable:
		return m.mapLuaTableToGoStruct(v, rv)
	case *lua.LUserData:
		return mapLuaUserDataToGoValue(v, rv)
	}
	return newTypeError(lv, rv)
}

func (m *Mapper) mapLuaTableToGoStruct(tbl *lua.LTable, rv reflect.Value) error {
	assert.True(tbl != nil)
	assert.True(rv.Kind() == reflect.Struct)
	rvType := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		fldVal := rv.Field(i)
		if !fldVal.CanSet() {
			continue // unexported field
		}

		field := rvType.Field(i)
		fieldName := getFieldName(field, m.TagName)
		if err := m.MapValue(tbl.RawGet(lua.LString(fieldName)), fldVal); err != nil {
			return fmt.Errorf("%s: %w", field.Name, err)
		}
	}
	return nil
}

// getFieldName get the struct field name.
func getFieldName(field reflect.StructField, tagName string) string {
	fieldName := field.Name
	if tagName == "" {
		return fieldName
	}

	tagValue := field.Tag.Get(tagName)
	tagSubValue := strings.SplitN(tagValue, ",", 2)[0]
	if tagSubValue != "" {
		return tagSubValue // use field name from tag value
	}
	return fieldName
}

func (m *Mapper) mapLuaTableToGoMap(tbl *lua.LTable, rv reflect.Value) error {
	assert.True(tbl != nil)
	assert.True(rv.Kind() == reflect.Map)
	mapType := rv.Type()
	keyType := mapType.Key()
	elemType := mapType.Elem()
	if rv.IsNil() || rv.Len() > 0 { // reset map
		rv.Set(reflect.MakeMap(mapType))
	}
	tbl.ForEach(func(lKey, lVal lua.LValue) {
		rvKeyPtr := reflect.New(keyType) // rvKeyPtr is a pointer to a new zero key
		rvKey := rvKeyPtr.Elem()
		if err := m.MapValue(lKey, rvKeyPtr.Elem()); err != nil {
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
