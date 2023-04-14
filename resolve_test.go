package ioc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTryResolveNamedContext(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(context.Background(), "foo", "bar")

		c := new(Container)
		var outFoo interface{}
		BindNamed(c, "fizz", func(c *Container) (int, error) {
			outFoo = c.Context().Value("foo")
			return 123, nil
		})

		out, err := TryResolveNamedContext[int](ctx, c, "fizz")
		assert.NoError(t, err)
		assert.Equal(t, 123, out)
		assert.Equal(t, "bar", outFoo)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		c := new(Container)
		exErr := errors.New("some error")
		BindNamed(c, "fizz", func(_ *Container) (int, error) {
			return 42, exErr
		})

		out, err := TryResolveNamedContext[int](context.Background(), c, "fizz")
		assert.Zero(t, out)
		assert.Equal(t, exErr, err)
	})

	t.Run("missing provider", func(t *testing.T) {
		t.Parallel()

		c := new(Container)

		out, err := TryResolveNamedContext[int](context.Background(), c, "fizz")
		assert.Zero(t, out)
		assert.ErrorAs(t, err, &MissingProviderError{})
	})

	t.Run("circular reference", func(t *testing.T) {
		t.Parallel()

		c := new(Container)
		BindNamed(c, "foo", func(c *Container) (int, error) {
			return TryResolveNamedContext[int](c.Context(), c, "bar")
		})
		BindNamed(c, "bar", func(c *Container) (int, error) {
			return TryResolveNamedContext[int](c.Context(), c, "foo")
		})

		out, err := TryResolveNamedContext[int](context.Background(), c, "foo")
		assert.Zero(t, out)
		assert.ErrorAs(t, err, &CircularDependencyError{})
	})
}

func TestTryResolve(t *testing.T) {
	t.Parallel()

	c := new(Container)
	Bind(c, Static(42))

	out, err := TryResolve[int](c)
	assert.Equal(t, 42, out)
	assert.NoError(t, err)
}

func TestTryResolveContext(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "foo", "bar")

	c := new(Container)
	var outCtx interface{}
	Bind(c, Infallible(func(c *Container) int {
		outCtx = c.Context().Value("foo")
		return 42
	}))

	out, err := TryResolveContext[int](ctx, c)
	assert.Equal(t, 42, out)
	assert.NoError(t, err)
	assert.Equal(t, "bar", outCtx)
}

func TestTryResolveNamed(t *testing.T) {
	t.Parallel()

	c := new(Container)
	BindNamed(c, "foo", Static(42))

	out, err := TryResolveNamed[int](c, "foo")
	assert.Equal(t, 42, out)
	assert.NoError(t, err)
}

func TestResolve(t *testing.T) {
	t.Parallel()

	c := new(Container)
	Bind(c, Static(42))

	out := Resolve[int](c)
	assert.Equal(t, 42, out)
	assert.Panics(t, func() { Resolve[string](c) })
}

func TestResolveContext(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "foo", "bar")

	c := new(Container)
	var outCtx interface{}
	Bind(c, Infallible(func(c *Container) int {
		outCtx = c.Context().Value("foo")
		return 42
	}))

	out := ResolveContext[int](ctx, c)
	assert.Equal(t, 42, out)
	assert.Equal(t, "bar", outCtx)

	assert.Panics(t, func() { ResolveContext[string](ctx, c) })
}

func TestResolveNamed(t *testing.T) {
	t.Parallel()

	c := new(Container)
	BindNamed(c, "foo", Static(42))

	out := ResolveNamed[int](c, "foo")
	assert.Equal(t, 42, out)
	assert.Panics(t, func() { ResolveNamed[string](c, "bar") })
}

func TestResolveNamedContext(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "foo", "bar")

	c := new(Container)
	var outFoo interface{}
	BindNamed(c, "fizz", func(c *Container) (int, error) {
		outFoo = c.Context().Value("foo")
		return 123, nil
	})

	out := ResolveNamedContext[int](ctx, c, "fizz")
	assert.Equal(t, 123, out)
	assert.Equal(t, "bar", outFoo)
	assert.Panics(t, func() { ResolveNamedContext[string](ctx, c, "bar") })
}
