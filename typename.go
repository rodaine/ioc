package ioc

import (
	"fmt"
	"reflect"
)

type typeName struct {
	Name string
	Type reflect.Type
}

func newTypeName[T any](name string) typeName {
	return typeName{
		Name: name,
		Type: reflect.TypeOf((*T)(nil)).Elem(),
	}
}

func (tn typeName) String() string {
	if tn.Name == anonymous {
		return tn.Type.String()
	}
	return fmt.Sprintf("%v:%s", tn.Type, tn.Name)
}

var _ fmt.Stringer = typeName{}
