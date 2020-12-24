package gluamapper

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/yuin/gopher-lua"
)

type testRole struct {
	Name string
}

type testPerson struct {
	Name      string
	Age       int
	WorkPlace string
	Role      []*testRole
}

func TestReflectValueSetZero(t *testing.T) {
	var rv reflect.Value
	assert := require.New(t)
	var n int
	rv = reflect.ValueOf(&n).Elem()
	rv.Set(reflect.Zero(rv.Type()))
	assert.Zero(n)

	var itf interface{}
	rv = reflect.ValueOf(&itf).Elem()
	rv.Set(reflect.Zero(rv.Type()))
	assert.Empty(itf)

	var arr [3]int
	rv = reflect.ValueOf(&arr).Elem()
	rv.Set(reflect.Zero(rv.Type()))
	assert.Zero(arr)
	var st struct{ a int }
	rv = reflect.ValueOf(&st).Elem()
	rv.Set(reflect.Zero(rv.Type()))
	assert.Zero(st)
	var p *int
	rv = reflect.ValueOf(&p).Elem()
	rv.Set(reflect.Zero(rv.Type()))
	assert.Zero(p)
}

func TestTagName(t *testing.T) {
	var err error
	assert := require.New(t)
	L := lua.NewState()
	err = L.DoString(`
		a = {aabbcc = 123, Abc = 456}
	`)
	assert.NoError(err)

	type A struct {
		Abc int `mytag:"aabbcc"`
	}
	m := NewMapperWithTagName("mytag")
	var output A
	err = m.Map(L.GetGlobal("a"), &output)
	assert.NoError(err)
	assert.Equal(A{Abc: 123}, output)

	type B struct {
		Bbb int `mytag:""`
	}
	err = L.DoString(`b = {Bbb =  123}`)
	assert.NoError(err)
	var goB B
	err = m.Map(L.GetGlobal("b"), &goB)
	assert.Equal(B{Bbb: 123}, goB)
}

func TestMapMap(t *testing.T) {
	var err error
	var output map[int]int
	assert := require.New(t)

	L := lua.NewState()
	err = L.DoString(`
		tbl = {abc = 123, [222]=222, [333]="333", [444]=444.4}
		arr = {1,2,3}
		n = 123
	`)
	assert.NoError(err)

	err = Map(L.GetGlobal("tbl"), &output)
	assert.NoError(err)
	assert.Equal(2, len(output))
	assert.Equal(222, output[222])
	assert.Equal(444, output[444]) // 444.4 -> 444

	err = Map(L.GetGlobal("arr"), &output)
	assert.NoError(err)
	assert.Equal(3, len(output))
	assert.Equal(map[int]int{1: 1, 2: 2, 3: 3}, output)

	err = Map(L.GetGlobal("n"), &output)
	assert.EqualError(err, "map[int]int expected but got Lua number")

	ud := L.NewUserData()
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Nil(output)
	ud.Value = 123
	err = Map(ud, &output)
	assert.EqualError(err, "map[int]int expected but got Lua user data of int")
	ud.Value = map[int]int{1: 1}
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(1, output[1])
	assert.Equal(1, len(output))
}

func TestMapPtr(t *testing.T) {
	var err error
	var output *int
	assert := require.New(t)

	err = Map(lua.LNumber(123), &output)
	assert.NoError(err)
	assert.Equal(123, *output)

	err = Map(lua.LNil, &output)
	assert.NoError(err)
	assert.Nil(output)

	err = Map(lua.LString("abc"), &output)
	assert.EqualError(err, "int expected but got Lua string")

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
	assert.EqualError(err, "*int expected but got Lua user data of *float64")
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
	output = []int{1}
	err = Map(tbl, &output)
	assert.NoError(err)
	assert.Equal([]int{1, 2, 3}, output)
	output = []int{4, 5, 6}
	err = Map(tbl, &output)
	assert.NoError(err)
	assert.Equal([]int{1, 2, 3}, output)
	output = []int{1, 2, 3, 4, 5, 6, 7}
	err = Map(tbl, &output)
	assert.NoError(err)
	assert.Equal([]int{1, 2, 3}, output)

	err = L.DoString(`t = {1,2,3,true}`)
	assert.NoError(err)
	tbl = L.GetGlobal("t")
	err = Map(tbl, &output)
	assert.EqualError(err, "slice[3]: int expected but got Lua boolean")

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
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal(goSlice, output)

	goArray := [3]int{1, 2, 3}
	ud.Value = goArray
	err = Map(ud, &output)
	assert.EqualError(err, "[]int expected but got Lua user data of [3]int")

	goFloatSlice := []float32{1, 2, 3}
	ud.Value = goFloatSlice
	err = Map(ud, &output)
	assert.EqualError(err, "[]int expected but got Lua user data of []float32")

	err = Map(lua.LTrue, &output)
	assert.EqualError(err, "[]int expected but got Lua boolean")
}

func TestMapStruct(t *testing.T) {
	var err error
	var output testPerson
	assert := require.New(t)

	err = Map(lua.LNil, &output)
	assert.NoError(err)
	err = Map(lua.LTrue, &output)
	assert.EqualError(err, "gluamapper.testPerson expected but got Lua boolean")

	L := lua.NewState()
	err = L.DoString(`
	    person = { Name = "Michel" }
		p2 = { Name = true }
	`)
	assert.NoError(err)
	err = Map(L.GetGlobal("person"), &output)
	assert.NoError(err)
	assert.Equal(testPerson{Name: "Michel"}, output)
	err = Map(L.GetGlobal("p2"), &output)
	assert.EqualError(err, "Name: string expected but got Lua boolean")

	ud := L.NewUserData()
	ud.Value = nil
	err = Map(ud, &output)
	assert.EqualError(err, "gluamapper.testPerson expected but got Lua user data of nil")

	ud.Value = testPerson{Name: "name"}
	output.Name = ""
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Equal("name", output.Name)

	ud.Value = &testPerson{}
	err = Map(ud, &output)
	assert.EqualError(err, "gluamapper.testPerson expected but got Lua user data of *gluamapper.testPerson")
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

	err = L.DoString(`
		t = {Aa = 123, bb = 456}
	`)
	assert.NoError(err)
	err = Map(L.GetGlobal("t"), &output)
	assert.NoError(err)
	assert.Equal(A{Aa: 123, bb: 456}, output)

	err = L.DoString(`
		t = {wall = 1234}
	`)
	assert.NoError(err)
	var tm time.Time // struct { wall uint64, ...}
	err = Map(L.GetGlobal("t"), &tm)
	assert.NoError(err)
	assert.Equal(time.Time{}, tm) // wall is unexported
}

func TestValueOfNil(t *testing.T) {
	assert := require.New(t)
	var x interface{}
	rv := reflect.ValueOf(x)
	assert.Equal(reflect.Invalid, rv.Kind())
}

func TestNilValue(t *testing.T) {
	assert := require.New(t)
	var p *int
	err := Map(lua.LTrue, p)
	assert.Equal(OutputValueIsNilError, err)
}

func TestMapValueIsNil(t *testing.T) {
	assert := require.New(t)
	var n *int
	err := Map(lua.LTrue, n)
	assert.EqualError(err, "output value is nil")
}

func TestMapArray(t *testing.T) {
	assert := require.New(t)
	L := lua.NewState()
	err := L.DoString(`tbl = {1,2,3,4,5.5,a=123}`)
	assert.NoError(err)
	tbl := L.GetGlobal("tbl")

	var a [3]int
	err = Map(tbl, &a)
	assert.NoError(err)
	assert.Equal([3]int{1, 2, 3}, a)
	var b [10]int
	err = Map(tbl, &b)
	assert.NoError(err)
	assert.Equal([10]int{1, 2, 3, 4, 5}, b)
	var c [2]bool
	err = Map(tbl, &c)
	assert.EqualError(err, "array[0]: bool expected but got Lua number")

	err = L.DoString(`t = 1234`)
	assert.NoError(err)
	err = Map(L.GetGlobal("t"), &a)
	assert.EqualError(err, "[3]int expected but got Lua number")

	ud := L.NewUserData()
	arr := [3]int{4, 5, 6}
	ud.Value = arr
	var d [3]int
	err = Map(ud, &d)
	assert.Equal(arr, d)
	err = Map(ud, &c)
	assert.EqualError(err, "[2]bool expected but got Lua user data of [3]int")
}
