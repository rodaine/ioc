package ioc

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCircularDependencyError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		err   CircularDependencyError
		exMsg string
	}{
		{
			name:  "empty",
			err:   CircularDependencyError{},
			exMsg: "circular dependency encountered",
		},
		{
			name:  "one",
			err:   CircularDependencyError{newTypeName[uint](anonymous)},
			exMsg: "circular dependency encountered resolving uint",
		},
		{
			name: "many",
			err: CircularDependencyError{
				newTypeName[int]("foo"),
				newTypeName[string]("bar"),
				newTypeName[bool]("baz"),
			},
			exMsg: "circular dependency encountered resolving int:foo:",
		},
	}

	for _, test := range tests {
		tc := test
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var msg string
			require.NotPanics(t, func() { msg = tc.err.Error() })
			assert.Contains(t, msg, tc.exMsg)
		})
	}
}

func TestMissingProviderError_Error(t *testing.T) {
	t.Parallel()

	err := MissingProviderError(newTypeName[io.Writer]("w"))
	assert.Contains(t, err.Error(), "missing provider for io.Writer:w")
}
