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

func ExampleMapper_Map_slice() {
	L := lua.NewState()
	if err := L.DoString(`a = {1, 2, 3.3}`); err != nil {
		panic(err)
	}

	var ints []int
	if err := NewMapper().Map(L.GetGlobal("a"), &ints); err != nil {
		panic(err)
	}
	fmt.Printf("ints: %v\n", ints)

	var floats []float32
	if err := NewMapper().Map(L.GetGlobal("a"), &floats); err != nil {
		panic(err)
	}
	fmt.Printf("floats: %v\n", floats)

	// Output:
	// ints: [1 2 3]
	// floats: [1 2 3.3]
}

func ExampleMapper_Map_map() {
	L := lua.NewState()
	if err := L.DoString(`a = {a=1, b=2, c=3.3, d=true, [123]=123}`); err != nil {
		panic(err)
	}

	var output map[string]int
	if err := NewMapper().Map(L.GetGlobal("a"), &output); err != nil {
		panic(err)
	}
	fmt.Printf("%v", output)
	// Output:
	// map[a:1 b:2 c:3]
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
