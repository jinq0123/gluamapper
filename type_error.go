package gluamapper

import (
	"fmt"
	"reflect"

	assert "github.com/arl/assertgo"
	"github.com/yuin/gopher-lua"
)

type TypeError struct {
	goType  reflect.Type
	luaType lua.LValueType

	// if luaType is LTUserData
	isLuaUserDataValueNil bool
	luaUserDataValueType  reflect.Type
}

func newTypeError(lv lua.LValue, rv reflect.Value) *TypeError {
	assert.True(rv.IsValid())
	goType := rv.Type()
	luaType := lv.Type()
	result := &TypeError{
		goType:  goType,
		luaType: luaType,
	}
	if luaType != lua.LTUserData {
		return result
	}

	ud := lv.(*lua.LUserData)
	if ud.Value == nil {
		result.isLuaUserDataValueNil = true
		return result
	}

	val := reflect.ValueOf(ud.Value)
	result.luaUserDataValueType = reflect.Indirect(val).Type()
	return result
}

func (t *TypeError) Error() string {
	if t.luaType != lua.LTUserData {
		return fmt.Sprintf("%s expected but got Lua %s", t.goType, t.luaType)
	}
	if t.isLuaUserDataValueNil {
		return fmt.Sprintf("%s expected but got Lua user data of nil", t.goType)
	}
	return fmt.Sprintf("%s expected but got Lua user data of %s", t.goType, t.luaUserDataValueType)
}
