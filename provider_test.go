package ioc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSingleton(t *testing.T) {
	t.Parallel()

	counter := 0
	counterProvider := func(_ Resolver) (int, error) {
		counter++
		return counter, nil
	}

	c := new(Container)
	Bind(c, counterProvider)
	BindNamed(c, "single", Singleton(counterProvider))

	r := c.Freeze()
	assert.Equal(t, 1, Resolve[int](r))
	assert.Equal(t, 2, Resolve[int](r))
	assert.Equal(t, 3, ResolveNamed[int](r, "single"))
	assert.Equal(t, 3, ResolveNamed[int](r, "single"))
	assert.Equal(t, 4, Resolve[int](r))
	assert.Equal(t, 3, ResolveNamed[int](r, "single"))
	assert.Equal(t, 4, counter)
}

func TestStatic(t *testing.T) {
	t.Parallel()

	c := new(Container)
	Bind(c, Static(123))

	r := c.Freeze()
	assert.Equal(t, 123, Resolve[int](r))
	assert.Equal(t, 123, Resolve[int](r))
	assert.Equal(t, 123, Resolve[int](r))
}

func TestInfallible(t *testing.T) {
	t.Parallel()

	fn := func(_ Resolver) string {
		return "foobar"
	}

	c := new(Container)
	Bind(c, Infallible(fn))

	r := c.Freeze()
	assert.Equal(t, "foobar", Resolve[string](r))
	assert.Equal(t, "foobar", Resolve[string](r))
	assert.Equal(t, "foobar", Resolve[string](r))
}
