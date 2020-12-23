package gluamapper

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yuin/gopher-lua"
)

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
