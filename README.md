# gluamapper: maps a GopherLua table to a Go struct

[![Build Status](https://travis-ci.org/jinq0123/gluamapper.svg)](https://travis-ci.org/jinq0123/gluamapper)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/jinq0123/gluamapper)](https://pkg.go.dev/github.com/jinq0123/gluamapper)

gluamapper provides an easy way to map GopherLua tables to Go structs.

gluamapper converts a GopherLua table to `map[string]interface{}`,
 and then converts it to a Go struct using [`mapstructure`](https://github.com/mitchellh/mapstructure/).
 
## Installation

```bash
go get github.com/jinq0123/gluamapper
```

## API
See [Go doc](https://pkg.go.dev/github.com/jinq0123/gluamapper).

## Usage

```go

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
    var person Person
    if err := gluamapper.Map(L.GetGlobal("person").(*lua.LTable), &person); err != nil {
        panic(err)
    }
    fmt.Printf("%s %d", person.Name, person.Age)
```

## License
MIT

## Author
* Yusuke Inuzuka
