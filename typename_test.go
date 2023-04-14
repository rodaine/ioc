package ioc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	renamedPkg "github.com/stretchr/testify/require" // intentionally renamed for testing purposes
)

func TestTypeName_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		tn typeName
		ex string
	}{
		{
			tn: newTypeName[int](anonymous),
			ex: "int",
		},
		{
			tn: newTypeName[int]("foo"),
			ex: "int:foo",
		},
		{
			tn: newTypeName[assert.TestingT]("bar"),
			ex: "assert.TestingT:bar",
		},
		{
			tn: newTypeName[renamedPkg.Assertions](anonymous),
			ex: "require.Assertions",
		},
	}

	for _, tc := range tests {
		test := tc
		t.Run(test.ex, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.ex, test.tn.String())
		})
	}
}
