package gluamapper

import (
	"fmt"
	"reflect"
)

type OutputIsNotAPointerError struct {
	outputValue reflect.Value
}

func (o *OutputIsNotAPointerError) Error() string {
	return fmt.Sprintf("output is not a pointer but a kind of %s", o.outputValue.Kind())
}
