package gluamapper

import (
	"fmt"

	"github.com/yuin/gopher-lua"
)

func ExampleMap() {
	type Role struct {
		Name string
	}

	type Person struct {
		Name      string
		Age       int
		WorkPlace string
		Role      []*Role
	}

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
		panic(err)
	}
	var person Person
	if err := Map(L.GetGlobal("person"), &person); err != nil {
		panic(err)
	}
	fmt.Printf("%s %d", person.Name, person.Age)
	// Output:
	// Michel 31
}
