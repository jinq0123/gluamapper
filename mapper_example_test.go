package gluamapper

import (
	"fmt"
	"reflect"

	"github.com/yuin/gopher-lua"
)

func ExampleMapper_Map() {
	L := lua.NewState()
	if err := L.DoString(`
      person = {
        name = "Michel",
        age  = 31,
      }
    `); err != nil {
		panic(err)
	}

	// Person struct with json tag
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var person Person
	mapper := NewMapperWithTagName("json") // must use json tag name
	if err := mapper.Map(L.GetGlobal("person"), &person); err != nil {
		panic(err)
	}
	fmt.Printf("%s %d", person.Name, person.Age)
	// Output:
	// Michel 31
}

func ExampleMapper_MapValue() {
	L := lua.NewState()
	if err := L.DoString(`
      person = {
        name = "Michel",
        age  = 31,
      }
    `); err != nil {
		panic(err)
	}

	// Person struct with json tag
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	var person Person
	mapper := NewMapperWithTagName("json") // must use json tag name
	rv := reflect.ValueOf(&person).Elem()
	if err := mapper.MapValue(L.GetGlobal("person"), rv); err != nil {
		panic(err)
	}
	fmt.Printf("%s %d", person.Name, person.Age)
	// Output:
	// Michel 31
}
