# gluamapper: maps a GopherLua table to a Go struct

[![Build Status](https://travis-ci.org/jinq0123/gluamapper.svg)](https://travis-ci.org/jinq0123/gluamapper)
[![codecov](https://codecov.io/gh/jinq0123/gluamapper/branch/master/graph/badge.svg?token=190O5EPVTH)](https://codecov.io/gh/jinq0123/gluamapper)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/jinq0123/gluamapper)](https://pkg.go.dev/github.com/jinq0123/gluamapper)

gluamapper provides an easy way to map GopherLua tables to Go structs.

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
* Jin Qing

## Differences from [yuin/gluamapper](https://github.com/yuin/gluamapper)

+ Speedup
	* Converts directly from lua table to go struct, while yuin/gluamapper
		converts the table to `map[string]interface{}`,
		and then converts it to a Go struct using [`mapstructure`](https://github.com/mitchellh/mapstructure/).
	* No "weak" conversions
		+ only lua number to int types allowed
	* Always ignores unused keys

+ New feature
	* Maps lua types other than table to go types

+ Bugfix
	* TODO: circular reference

## TODO
* handle embedded fields
