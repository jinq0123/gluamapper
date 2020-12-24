package gluamapper

import (
	"reflect"

	assert "github.com/arl/assertgo"
	"github.com/yuin/gopher-lua"
)

// Map maps the Lua value to the Go value pointed by output.
// If output is not a pointer, Map returns OutputIsNotAPointerError.
// If output is nil, Map returns OutputValueIsNilError.
// If the Lua value is nil, the Go value will be set to its zero value.
// Is the Lua value is a *lua.LUserData, the Go value will be set to the value of the LUserData
// if they are the same type, or TypeError will be returned if they are not the same type.
//
// Map will allocate maps, slices, and pointers as necessary,
// with the following additional rules:
//
// To map Lua value into a pointer, Map allocates a new value.
//
// To map Lua table into a struct, Map matches incoming Lua table
// keys to the struct field name or its tag.
// Lua table keys which don't have a corresponding struct field are ignored.
//
// To map Lua value into an interface value,
// Map stores one of these in the interface value:
//
//	bool, for Lua booleans
//	float64, for Lua numbers
//	string, for Lua strings
//	[]interface{}, for Lua arrays
//	map[string]interface{}, for Lua tables
//	nil for Lua nil
//
// To map a Lua array into a slice, Map resets the slice to a new one
// if the capacity is not large enough, or set the len of the slice.
//
// To map a Lua array into a Go array, Map maps Lua array elements into
// corresponding Go array elements.
// If the Go array is smaller than the Lua array,
// the additional Lua array elements are discarded.
// If the Lua array is smaller than the Go array,
// the additional Go array elements are set to zero values.
//
// To map a Lua table into a map, Map first allocates a map to use
// if the old map is nil or not empty.
// Map then stores key-value pairs from the Lua table into the map.
// The Lua table's key-values are ignored
// if the Lua key can not be mapped into a Go key
// or the Lua value can not be mapped into a Go value
//
// If tag name is needed, please use NewMapperWithTagName(tagName).Map(...)
func Map(lv lua.LValue, output interface{}) error {
	return NewMapper().Map(lv, output)
}

func mapBool(lv lua.LValue, rv reflect.Value) error {
	assert.True(rv.Kind() == reflect.Bool)
	if b, ok := toBool(lv); ok {
		rv.SetBool(b)
		return nil
	}
	return newTypeError(lv, rv)
}

func toBool(lv lua.LValue) (result bool, ok bool) {
	switch v := lv.(type) {
	case lua.LBool:
		return bool(v), true
	case *lua.LUserData:
		result, ok = v.Value.(bool)
		return result, ok
	}
	return false, false
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

// Always returns nil
func mapInterface(lv lua.LValue, rv reflect.Value) error {
	assert.True(lv != lua.LNil)
	assert.True(rv.Kind() == reflect.Interface)
	itf := toInterface(lv)
	if itf != nil {
		rv.Set(reflect.ValueOf(itf))
		return nil
	}

	// can not call of reflect.Value.Set on zero Value
	var i interface{} // nil
	ri := reflect.ValueOf(&i).Elem()
	rv.Set(ri) // Set to nil
	return nil
}

func toInterface(lv lua.LValue) interface{} {
	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LNumber:
		return float64(v)
	case lua.LString:
		return string(v)
	case *lua.LFunction:
		return v // keep as *LFunction
	case *lua.LUserData:
		return v.Value // may be nil
	// case *lua.LTThread: no such type
	case *lua.LTable:
		return luaTableToGoInterface(v)
	case lua.LChannel:
		return (chan lua.LValue)(v)
	default:
		return v // keep as v
	}
	return nil
}

func mapString(lv lua.LValue, rv reflect.Value) error {
	assert.True(lv != lua.LNil)
	assert.True(rv.Kind() == reflect.String)
	if s, ok := toString(lv); ok {
		rv.SetString(s)
		return nil
	}
	return newTypeError(lv, rv)
}

func toString(lv lua.LValue) (result string, ok bool) {
	switch v := lv.(type) {
	case lua.LString:
		return string(v), true
	case *lua.LUserData:
		result, ok := v.Value.(string)
		return result, ok
	}
	return "", false
}

func luaTableToGoMap(tbl *lua.LTable) map[string]interface{} {
	mp := make(map[string]interface{})
	tbl.ForEach(func(lKey, lVal lua.LValue) {
		if key, ok := toString(lKey); ok {
			mp[key] = toInterface(lVal)
		}
	})
	return mp
}

func luaTableToGoInterface(tbl *lua.LTable) interface{} {
	assert.True(tbl != nil)
	maxn := tbl.MaxN()
	if maxn == 0 { // table -> map[string]interface{}
		return luaTableToGoMap(tbl) // Only support string key
	}

	// else: array -> []interface{}
	slc := make([]interface{}, maxn, maxn)
	for i := 0; i < maxn; i++ {
		slc[i] = toInterface(tbl.RawGetInt(i + 1))
	}
	return slc
}

func mapLuaUserDataToGoValue(ud *lua.LUserData, rv reflect.Value) error {
	assert.True(rv.IsValid()) // rv.Kind() != Invalid
	udValue := ud.Value
	if udValue == nil {
		if canBeNil(rv) {
			rv.Set(reflect.Zero(rv.Type()))
			return nil
		}
		return &TypeError{
			goType:                rv.Type(),
			luaType:               lua.LTUserData,
			isLuaUserDataValueNil: true,
		}
	}

	// must be the same type
	udValType := reflect.TypeOf(udValue)
	if udValType == rv.Type() {
		rv.Set(reflect.ValueOf(udValue))
		return nil
	}

	return &TypeError{
		goType:               rv.Type(),
		luaType:              lua.LTUserData,
		luaUserDataValueType: udValType,
	}
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
