package gluamapper

import (
	"reflect"

	assert "github.com/arl/assertgo"
	"github.com/yuin/gopher-lua"
)

// Map maps the lua value to the given go pointer with default options.
// Please reset output before Map()
func Map(lv lua.LValue, output interface{}) error {
	return NewMapper().Map(lv, output)
}

func mapBool(lv lua.LValue, rv reflect.Value) error {
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

func mapInt(lv lua.LValue, rv reflect.Value) error {
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

func mapInt8(lv lua.LValue, rv reflect.Value) error {
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

func mapInt16(lv lua.LValue, rv reflect.Value) error {
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

func mapInt32(lv lua.LValue, rv reflect.Value) error {
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

func mapInt64(lv lua.LValue, rv reflect.Value) error {
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

func mapUint(lv lua.LValue, rv reflect.Value) error {
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

func mapUint8(lv lua.LValue, rv reflect.Value) error {
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

func mapUint16(lv lua.LValue, rv reflect.Value) error {
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

func mapUint32(lv lua.LValue, rv reflect.Value) error {
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

func mapUint64(lv lua.LValue, rv reflect.Value) error {
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

func mapFloat32(lv lua.LValue, rv reflect.Value) error {
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

func mapFloat64(lv lua.LValue, rv reflect.Value) error {
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

// TODO: need to return error?
func mapInterface(lv lua.LValue, rv reflect.Value) error {
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
		rv.Set(reflect.ValueOf(v)) // keep as *LFunction
	case *lua.LUserData:
		return mapLuaUserDataToGoInterface(v, rv)
	// case *lua.LTThread: no such type
	case *lua.LTable:
		return mapLuaTableToGoInterface(v, rv)
	case lua.LChannel:
		rv.Set(reflect.ValueOf((chan lua.LValue)(v)))
	default:
		rv.Set(reflect.ValueOf(v)) // keep as v
	}
	return nil
}

func mapString(lv lua.LValue, rv reflect.Value) error {
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

func mapLuaTableToGoInterface(tbl *lua.LTable, rv reflect.Value) error {
	assert.True(tbl != nil)
	assert.True(rv.Kind() == reflect.Interface)
	maxn := tbl.MaxN()
	if maxn == 0 { // table -> map[interface{}]interface{}
		// TODO: extract tblToMap()
		mp := make(map[interface{}]interface{})
		tbl.ForEach(func(lKey, lVal lua.LValue) {
			var key interface{}
			if err := mapInterface(lKey, reflect.ValueOf(&key).Elem()); err != nil {
				return // skip field if error
			}
			var val interface{}
			if err := mapInterface(lVal, reflect.ValueOf(&val).Elem()); err != nil {
				return // skip field if error
			}
			mp[key] = val
		})
		rv.Set(reflect.ValueOf(mp))
		return nil
	}

	// else: array -> []interface{}
	slc := make([]interface{}, maxn, maxn)
	rvSlc := reflect.ValueOf(slc)
	for i := 0; i < maxn; i++ {
		if err := mapInterface(tbl.RawGetInt(i+1), rvSlc.Index(i)); err != nil {
			return err
		}
	}
	rv.Set(rvSlc)
	return nil
}

// TODO: delete return
// change to luaUserDataToGoInterface(ud) interface{}
func mapLuaUserDataToGoInterface(ud *lua.LUserData, rv reflect.Value) error {
	assert.True(ud != nil)
	assert.True(rv.Kind() == reflect.Interface)
	udValue := ud.Value
	// can not call of reflect.Value.Set on zero Value
	if udValue != nil {
		rv.Set(reflect.ValueOf(udValue))
		return nil
	}

	var i interface{} // nil
	ri := reflect.ValueOf(&i).Elem()
	rv.Set(ri) // Set to nil
	return nil
}

func (m *Mapper) mapLuaTableToGoMap(tbl *lua.LTable, rv reflect.Value) error {
	assert.True(tbl != nil)
	assert.True(rv.Kind() == reflect.Map)
	mapType := rv.Type()
	keyType := mapType.Key()
	elemType := mapType.Elem()
	if rv.IsNil() {
		rv.Set(reflect.MakeMap(mapType))
	}
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
