package gluamapper

import (
	"reflect"
	"testing"

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

	output = nil
	err = Map(L.GetGlobal("arr"), &output)
	assert.NoError(err)
	assert.Equal(3, len(output))
	assert.Equal(map[int]int{1: 1, 2: 2, 3: 3}, output)

	err = Map(L.GetGlobal("n"), &output)
	assert.EqualError(err, "map[int]int expected but got lua number")

	ud := L.NewUserData()
	err = Map(ud, &output)
	assert.NoError(err)
	assert.Nil(output)
	ud.Value = 123
	err = Map(ud, &output)
	assert.EqualError(err, "map[int]int expected but got lua user data of int")
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

func TestValueOfNil(t *testing.T) {
	assert := require.New(t)
	var x interface{}
	rv := reflect.ValueOf(x)
	assert.Equal(reflect.Invalid, rv.Kind())
}
