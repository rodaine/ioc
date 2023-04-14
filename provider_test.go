package ioc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleton(t *testing.T) {
	t.Parallel()

	counter := 0
	counterProvider := func(_ *Container) (int, error) {
		counter++
		return counter, nil
	}

	c := new(Container)
	Bind(c, counterProvider)
	BindNamed(c, "single", Singleton(counterProvider))

	assert.Equal(t, 1, Resolve[int](c))
	assert.Equal(t, 2, Resolve[int](c))
	assert.Equal(t, 3, ResolveNamed[int](c, "single"))
	assert.Equal(t, 3, ResolveNamed[int](c, "single"))
	assert.Equal(t, 4, Resolve[int](c))
	assert.Equal(t, 3, ResolveNamed[int](c, "single"))
	assert.Equal(t, 4, counter)
}

func TestStatic(t *testing.T) {
	t.Parallel()

	c := new(Container)
	Bind(c, Static(123))

	assert.Equal(t, 123, Resolve[int](c))
	assert.Equal(t, 123, Resolve[int](c))
	assert.Equal(t, 123, Resolve[int](c))
}

func TestInfallible(t *testing.T) {
	t.Parallel()

	fn := func(_ *Container) string {
		return "foobar"
	}

	c := new(Container)
	Bind(c, Infallible(fn))

	assert.Equal(t, "foobar", Resolve[string](c))
	assert.Equal(t, "foobar", Resolve[string](c))
	assert.Equal(t, "foobar", Resolve[string](c))
}
