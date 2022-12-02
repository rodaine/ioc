package ioc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	// intentionally renamed for testing purposes:
	renamedPkg "github.com/stretchr/testify/require"
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

	for _, test := range tests {
		t.Run(test.ex, func(t *testing.T) {
			assert.Equal(t, test.ex, test.tn.String())
		})
	}
}
