package gluamapper

import (
	"testing"

	yuin "github.com/yuin/gluamapper"
	"github.com/yuin/gopher-lua"
)

func BenchmarkJinq0123Map(b *testing.B) {
	tbl := getLuaPersonTable()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var person testPerson
		if err := Map(tbl, &person); err != nil {
			panic(err)
		}
	}
}

func BenchmarkYuinMap(b *testing.B) {
	tbl := getLuaPersonTable()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var person testPerson
		if err := yuin.Map(tbl, &person); err != nil {
			panic(err)
		}
	}
}

func getLuaPersonTable() *lua.LTable {
	L := lua.NewState()
	if err := L.DoString(`
    person = {
      name = "Michel",
      age  = "31", -- weakly input
      work_place = "San Jose",
      role = {
        {
          name = "Administrator"
        },
        {
          name = "Operator"
        }
      }
    }
	`); err != nil {
		panic(err)
	}
	return L.GetGlobal("person").(*lua.LTable)
}
