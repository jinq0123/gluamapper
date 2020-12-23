package gluamapper

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
)

func errorIfNotEqual(t *testing.T, v1, v2 interface{}) {
	if v1 != v2 {
		_, file, line, _ := runtime.Caller(1)
		t.Errorf("%v line %v: '%v' expected, but got '%v'", filepath.Base(file), line, v1, v2)
	}
}

type testRole struct {
	Name string
}

type testPerson struct {
	Name      string
	Age       int
	WorkPlace string
	Role      []*testRole
}

type testStruct struct {
	Nil    interface{}
	Bool   bool
	String string
	Number int
	Func   interface{}
}

func TestMap(t *testing.T) {
	L := lua.NewState()
	if err := L.DoString(`
    person = {
      Name = "Michel",
      Age  = 31,
      WorkPlace = "San Jose",
      Role = {
        {
          Name = "Administrator"
        },
        {
          Name = "Operator"
        }
      }
    }
	`); err != nil {
		t.Error(err)
	}
	var person testPerson
	if err := Map(L.GetGlobal("person").(*lua.LTable), &person); err != nil {
		t.Error(err)
	}
	errorIfNotEqual(t, "Michel", person.Name)
	errorIfNotEqual(t, 31, person.Age)
	errorIfNotEqual(t, "San Jose", person.WorkPlace)
	errorIfNotEqual(t, 2, len(person.Role))
	errorIfNotEqual(t, "Administrator", person.Role[0].Name)
	errorIfNotEqual(t, "Operator", person.Role[1].Name)
}

func TestTypes(t *testing.T) {
	L := lua.NewState()
	if err := L.DoString(`
    tbl = {
      ["Nil"] = nil,
      ["Bool"] = true,
      ["String"] = "string",
      ["Number"] = 10,
      ["Func"] = function() end
    }
	`); err != nil {
		t.Error(err)
	}
	var stct testStruct

	if err := Map(L.GetGlobal("tbl"), &stct); err != nil {
		t.Error(err)
	}
	errorIfNotEqual(t, nil, stct.Nil)
	errorIfNotEqual(t, true, stct.Bool)
	errorIfNotEqual(t, "string", stct.String)
	errorIfNotEqual(t, 10, stct.Number)
}

func TestError(t *testing.T) {
	assert := require.New(t)
	L := lua.NewState()
	tbl := L.NewTable()
	L.SetField(tbl, "key", lua.LString("value"))
	err := Map(tbl, 1)
	assert.EqualError(err, "output must be a pointer")
}

func TestMapBool(t *testing.T) {
	var output bool
	var err error
	assert := require.New(t)
	L := lua.NewState()

	L.SetGlobal("goTrue", lua.LTrue)
	err = Map(L.GetGlobal("goTrue"), &output)
	assert.NoError(err)
	assert.Equal(true, output)
	L.SetGlobal("goFalse", lua.LFalse)
	err = Map(L.GetGlobal("goFalse"), &output)
	assert.NoError(err)
	assert.Equal(false, output)
	L.SetGlobal("goInt", lua.LNumber(12345))
	err = Map(L.GetGlobal("goInt"), &output)
	assert.EqualError(err, "bool expected but got lua number")
	L.SetGlobal("goNil", lua.LNil)
	err = Map(L.GetGlobal("goNil"), &output)
	assert.NoError(err)
	assert.Equal(false, output)
	goSt := struct{ a int }{a: 1234}
	L.SetGlobal("goSt", luar.New(L, &goSt))
	err = Map(L.GetGlobal("goSt"), &output)
	assert.EqualError(err, "bool expected but got lua user data of struct { a int }")
	ud := L.NewUserData()
	ud.Value = true
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(true, output)
	ud.Value = 1234
	err = Map(ud, &output)
	assert.EqualError(err, "bool expected but got lua user data of int")

	err = L.DoString(`
		luaTrue = true
		luaFalse = false
		luaInt = 123
		luaNil = nil
	`)
	assert.NoError(err)
	err = Map(L.GetGlobal("luaTrue"), &output)
	assert.NoError(err)
	assert.True(true, output)
	err = Map(L.GetGlobal("luaFalse"), &output)
	assert.NoError(err)
	assert.Equal(false, output)
	err = Map(L.GetGlobal("luaInt"), &output)
	assert.EqualError(err, "bool expected but got lua number")
	err = Map(L.GetGlobal("luaNil"), &output)
	assert.NoError(err)
	assert.Equal(false, output)
	err = Map(L.GetGlobal("luaNoSuchValue"), &output)
	assert.NoError(err)
	assert.Equal(false, output)
}

func TestMapInt(t *testing.T) {
	var err error
	var output int
	assert := require.New(t)
	L := lua.NewState()

	err = Map(lua.LNil, &output)
	assert.NoError(err)
	assert.Equal(int(0), output)
	err = Map(lua.LNumber(1234), &output)
	assert.NoError(err)
	assert.Equal(int(1234), output)
	err = Map(lua.LTrue, &output)
	assert.Error(err)
	assert.EqualError(err, "int expected but got lua boolean")
	ud := L.NewUserData()
	ud.Value = 1234
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(int(1234), output)
	ud.Value = "abce"
	err = Map(ud, &output)
	assert.EqualError(err, "int expected but got lua user data of string")
}

func TestMapInt8(t *testing.T) {
	var err error
	var output int8
	assert := require.New(t)
	L := lua.NewState()

	err = Map(lua.LNil, &output)
	assert.NoError(err)
	assert.Equal(int8(0), output)
	err = Map(lua.LNumber(12), &output)
	assert.NoError(err)
	assert.Equal(int8(12), output)
	err = Map(lua.LTrue, &output)
	assert.Error(err)
	assert.EqualError(err, "int8 expected but got lua boolean")
	ud := L.NewUserData()
	ud.Value = int8(12)
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(int8(12), output)
	ud.Value = "abce"
	err = Map(ud, &output)
	assert.EqualError(err, "int8 expected but got lua user data of string")
}

func TestMapInt16(t *testing.T) {
	var err error
	var output int16
	assert := require.New(t)
	L := lua.NewState()

	err = Map(lua.LNil, &output)
	assert.NoError(err)
	assert.Equal(int16(0), output)
	err = Map(lua.LNumber(1234), &output)
	assert.NoError(err)
	assert.Equal(int16(1234), output)
	err = Map(lua.LTrue, &output)
	assert.Error(err)
	assert.EqualError(err, "int16 expected but got lua boolean")
	ud := L.NewUserData()
	ud.Value = int16(1234)
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(int16(1234), output)
	ud.Value = "abce"
	err = Map(ud, &output)
	assert.EqualError(err, "int16 expected but got lua user data of string")
}

func TestMapInt32(t *testing.T) {
	var err error
	var output int32
	assert := require.New(t)
	L := lua.NewState()

	err = Map(lua.LNil, &output)
	assert.NoError(err)
	assert.Equal(int32(0), output)
	err = Map(lua.LNumber(1234), &output)
	assert.NoError(err)
	assert.Equal(int32(1234), output)
	err = Map(lua.LTrue, &output)
	assert.Error(err)
	assert.EqualError(err, "int32 expected but got lua boolean")
	ud := L.NewUserData()
	ud.Value = int32(1234)
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(int32(1234), output)
	ud.Value = "abce"
	err = Map(ud, &output)
	assert.EqualError(err, "int32 expected but got lua user data of string")
}

func TestMapInt64(t *testing.T) {
	var err error
	var output int64
	assert := require.New(t)
	L := lua.NewState()

	err = Map(lua.LNil, &output)
	assert.NoError(err)
	assert.Equal(int64(0), output)
	err = Map(lua.LNumber(1234), &output)
	assert.NoError(err)
	assert.Equal(int64(1234), output)
	err = Map(lua.LTrue, &output)
	assert.Error(err)
	assert.EqualError(err, "int64 expected but got lua boolean")
	ud := L.NewUserData()
	ud.Value = int64(1234)
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(int64(1234), output)
	ud.Value = "abce"
	err = Map(ud, &output)
	assert.EqualError(err, "int64 expected but got lua user data of string")
}

func TestMapUint(t *testing.T) {
	var err error
	var output uint
	assert := require.New(t)
	L := lua.NewState()

	err = Map(lua.LNil, &output)
	assert.NoError(err)
	assert.Equal(uint(0), output)
	err = Map(lua.LNumber(1234), &output)
	assert.NoError(err)
	assert.Equal(uint(1234), output)
	err = Map(lua.LTrue, &output)
	assert.Error(err)
	assert.EqualError(err, "uint expected but got lua boolean")
	ud := L.NewUserData()
	ud.Value = uint(1234)
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(uint(1234), output)
	ud.Value = "abce"
	err = Map(ud, &output)
	assert.EqualError(err, "uint expected but got lua user data of string")
}

func TestMapUint8(t *testing.T) {
	var err error
	var output uint8
	assert := require.New(t)
	L := lua.NewState()

	err = Map(lua.LNil, &output)
	assert.NoError(err)
	assert.Equal(uint8(0), output)
	err = Map(lua.LNumber(12), &output)
	assert.NoError(err)
	assert.Equal(uint8(12), output)
	err = Map(lua.LTrue, &output)
	assert.Error(err)
	assert.EqualError(err, "uint8 expected but got lua boolean")
	ud := L.NewUserData()
	ud.Value = uint8(12)
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(uint8(12), output)
	ud.Value = "abce"
	err = Map(ud, &output)
	assert.EqualError(err, "uint8 expected but got lua user data of string")
}

func TestMapUint16(t *testing.T) {
	var err error
	var output uint16
	assert := require.New(t)
	L := lua.NewState()

	err = Map(lua.LNil, &output)
	assert.NoError(err)
	assert.Equal(uint16(0), output)
	err = Map(lua.LNumber(1234), &output)
	assert.NoError(err)
	assert.Equal(uint16(1234), output)
	err = Map(lua.LTrue, &output)
	assert.Error(err)
	assert.EqualError(err, "uint16 expected but got lua boolean")
	ud := L.NewUserData()
	ud.Value = uint16(1234)
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(uint16(1234), output)
	ud.Value = "abce"
	err = Map(ud, &output)
	assert.EqualError(err, "uint16 expected but got lua user data of string")
}

func TestMapUint32(t *testing.T) {
	var err error
	var output uint32
	assert := require.New(t)
	L := lua.NewState()

	err = Map(lua.LNil, &output)
	assert.NoError(err)
	assert.Equal(uint32(0), output)
	err = Map(lua.LNumber(1234), &output)
	assert.NoError(err)
	assert.Equal(uint32(1234), output)
	err = Map(lua.LTrue, &output)
	assert.Error(err)
	assert.EqualError(err, "uint32 expected but got lua boolean")
	ud := L.NewUserData()
	ud.Value = uint32(1234)
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(uint32(1234), output)
	ud.Value = "abce"
	err = Map(ud, &output)
	assert.EqualError(err, "uint32 expected but got lua user data of string")
}

func TestMapUint64(t *testing.T) {
	var err error
	var output uint64
	assert := require.New(t)
	L := lua.NewState()

	err = Map(lua.LNil, &output)
	assert.NoError(err)
	assert.Equal(uint64(0), output)
	err = Map(lua.LNumber(1234), &output)
	assert.NoError(err)
	assert.Equal(uint64(1234), output)
	err = Map(lua.LTrue, &output)
	assert.Error(err)
	assert.EqualError(err, "uint64 expected but got lua boolean")
	ud := L.NewUserData()
	ud.Value = uint64(1234)
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(uint64(1234), output)
	ud.Value = "abce"
	err = Map(ud, &output)
	assert.EqualError(err, "uint64 expected but got lua user data of string")
}

func TestMapFloat32(t *testing.T) {
	var err error
	var output float32
	assert := require.New(t)
	L := lua.NewState()

	err = Map(lua.LNil, &output)
	assert.NoError(err)
	assert.Equal(float32(0), output)
	err = Map(lua.LNumber(1234), &output)
	assert.NoError(err)
	assert.Equal(float32(1234), output)
	err = Map(lua.LTrue, &output)
	assert.Error(err)
	assert.EqualError(err, "float32 expected but got lua boolean")
	ud := L.NewUserData()
	ud.Value = float32(1234)
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(float32(1234), output)
	ud.Value = "abce"
	err = Map(ud, &output)
	assert.EqualError(err, "float32 expected but got lua user data of string")
}

func TestMapFloat64(t *testing.T) {
	var err error
	var output float64
	assert := require.New(t)
	L := lua.NewState()

	err = Map(lua.LNil, &output)
	assert.NoError(err)
	assert.Equal(float64(0), output)
	err = Map(lua.LNumber(1234), &output)
	assert.NoError(err)
	assert.Equal(float64(1234), output)
	err = Map(lua.LTrue, &output)
	assert.Error(err)
	assert.EqualError(err, "float64 expected but got lua boolean")
	ud := L.NewUserData()
	ud.Value = float64(1234)
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(float64(1234), output)
	ud.Value = "abce"
	err = Map(ud, &output)
	assert.EqualError(err, "float64 expected but got lua user data of string")
}

func TestMapSlice(t *testing.T) {
	var err error
	var output []int
	assert := require.New(t)
	L := lua.NewState()

	err = L.DoString(`t = {1,2,3}`)
	assert.NoError(err)
	tbl := L.GetGlobal("t")
	err = Map(tbl, &output)
	assert.NoError(err)
	assert.Equal([]int{1, 2, 3}, output)
	output = nil
	err = Map(lua.LNil, &output)
	assert.NoError(err)
	assert.Nil(output)

	ud := L.NewUserData()
	output = []int{1, 2, 3}
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Nil(output)

	goSlice := []int{1, 2, 3}
	ud.Value = goSlice
	output = nil
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(goSlice, output)

	goArray := [3]int{1, 2, 3}
	ud.Value = goArray
	output = nil
	err = Map(ud, &output)
	assert.EqualError(err, "[]int expected but got lua user data of [3]int")

	goFloatSlice := []float32{1, 2, 3}
	ud.Value = goFloatSlice
	output = nil
	err = Map(ud, &output)
	assert.EqualError(err, "[]int expected but got lua user data of []float32")

	err = Map(lua.LTrue, &output)
	assert.EqualError(err, "[]int expected but got lua boolean")
}

func TestMapString(t *testing.T) {
	var err error
	var output string
	assert := require.New(t)
	L := lua.NewState()

	err = Map(lua.LString("abc"), &output)
	assert.NoError(err)
	assert.Equal("abc", output)
	output = ""

	err = Map(lua.LNil, &output)
	assert.NoError(err)
	assert.Equal("", output)

	ud := L.NewUserData()
	ud.Value = "abc"
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal("abc", output)
	ud.Value = 123
	err = Map(ud, &output)
	assert.EqualError(err, "string expected but got lua user data of int")
	ud.Value = nil
	err = Map(ud, &output)
	assert.EqualError(err, "string expected but got lua user data of nil")
}

func TestMapStruct(t *testing.T) {
	var err error
	var output testPerson
	assert := require.New(t)

	err = Map(lua.LNil, &output)
	assert.NoError(err)

	L := lua.NewState()
	ud := L.NewUserData()
	ud.Value = nil
	err = Map(ud, &output)
	assert.EqualError(err, "gluamapper.testPerson expected but got lua user data of nil")

	ud.Value = testPerson{Name: "name"}
	output.Name = ""
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal("name", output.Name)

	ud.Value = &testPerson{}
	err = Map(ud, &output)
	assert.EqualError(err, "gluamapper.testPerson expected but got lua user data of *gluamapper.testPerson")
}

func TestMapStructWithUnexportedField(t *testing.T) {
	var err error
	assert := require.New(t)

	type A struct {
		Aa int
		bb int
	}
	a := A{Aa: 123, bb: 456}
	L := lua.NewState()
	ud := L.NewUserData()
	ud.Value = a

	var output A
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(123, output.Aa)
	assert.Equal(456, output.bb) // unexported field
}

func TestMapInterface(t *testing.T) {
	var err error
	var output interface{}
	assert := require.New(t)

	L := lua.NewState()
	err = L.DoString(`
		n = 234
		b = true
		s = "abc"
		t = {a=1234}
		f = function() end
	`)
	assert.NoError(err)
	err = Map(L.GetGlobal("n"), &output)
	assert.NoError(err)
	assert.Equal(float64(234), output)
	err = Map(L.GetGlobal("b"), &output)
	assert.NoError(err)
	assert.Equal(true, output)
	err = Map(L.GetGlobal("s"), &output)
	assert.NoError(err)
	assert.Equal("abc", output)
	err = Map(L.GetGlobal("t"), &output)
	assert.NoError(err)
	assert.Equal(float64(1234), output.(map[interface{}]interface{})["a"])
	output = nil
	err = Map(L.GetGlobal("f"), &output)
	assert.NoError(err)
	assert.NotNil(output)

	ch := make(chan lua.LValue)
	L.SetGlobal("ch", lua.LChannel(ch))
	err = Map(L.GetGlobal("ch"), &output)
	assert.NoError(err)
	assert.Equal(output, ch)

	ud := L.NewUserData()
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Nil(output)
	ud.Value = 123
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(123, output)
	ud.Value = nil
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Nil(output)
	arr := [12]int{1, 2, 3}
	ud.Value = arr
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(arr, output)
}

func TestMapMap(t *testing.T) {
	var err error
	var output map[int]int
	assert := require.New(t)

	L := lua.NewState()
	err = L.DoString(`
		tbl = {abc = 123, [222]=222, [333]="333", [444]=444.4}
		arr = {1,2,3}
	`)
	assert.NoError(err)

	err = Map(L.GetGlobal("tbl"), &output)
	assert.NoError(err)
	assert.Equal(2, len(output))
	assert.Equal(222, output[222])
	assert.Equal(444, output[444]) // 444.4 -> 444

	output = nil
	err = Map(L.GetGlobal("arr"), &output)
	assert.NoError(err)
	assert.Equal(3, len(output))
	assert.Equal(map[int]int{1: 1, 2: 2, 3: 3}, output)
	// XXX
}

func TestMapPtr(t *testing.T) {
	var err error
	var output *int
	assert := require.New(t)

	err = Map(lua.LNumber(123), &output)
	assert.NoError(err)
	assert.Equal(123, *output)

	L := lua.NewState()
	ud := L.NewUserData()
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Nil(output)

	n := 123
	ud.Value = &n
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(&n, output)

	f := 123.0
	ud.Value = &f
	err = Map(ud, &output)
	assert.EqualError(err, "*int expected but got lua user data of *float64")
}
