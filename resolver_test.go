package ioc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolver_Context(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	r := Resolver{}
	assert.Equal(t, ctx, r.Context())

	r.ctx = context.WithValue(ctx, "foo", "bar")
	assert.Equal(t, r.ctx, r.Context())
}

func TestTryResolveNamedContext(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.WithValue(context.Background(), "foo", "bar")

		c := new(Container)
		var outFoo interface{}
		BindNamed(c, "fizz", func(r Resolver) (int, error) {
			outFoo = r.Context().Value("foo")
			return 123, nil
		})
		r := c.Freeze()

		out, err := TryResolveNamedContext[int](ctx, r, "fizz")
		assert.NoError(t, err)
		assert.Equal(t, 123, out)
		assert.Equal(t, "bar", outFoo)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		c := new(Container)
		exErr := errors.New("some error")
		BindNamed(c, "fizz", func(r Resolver) (int, error) {
			return 42, exErr
		})
		r := c.Freeze()

		out, err := TryResolveNamedContext[int](context.Background(), r, "fizz")
		assert.Zero(t, out)
		assert.Equal(t, exErr, err)
	})

	t.Run("missing provider", func(t *testing.T) {
		t.Parallel()

		c := new(Container)
		r := c.Freeze()

		out, err := TryResolveNamedContext[int](context.Background(), r, "fizz")
		assert.Zero(t, out)
		assert.ErrorAs(t, err, &MissingProviderError{})
	})

	t.Run("circular reference", func(t *testing.T) {
		t.Parallel()

		c := new(Container)
		BindNamed(c, "foo", func(r Resolver) (int, error) {
			return TryResolveNamedContext[int](r.Context(), r, "bar")
		})
		BindNamed(c, "bar", func(r Resolver) (int, error) {
			return TryResolveNamedContext[int](r.Context(), r, "foo")
		})
		r := c.Freeze()

		out, err := TryResolveNamedContext[int](context.Background(), r, "foo")
		assert.Zero(t, out)
		assert.ErrorAs(t, err, &CircularDependencyError{})
	})
}

func TestTryResolve(t *testing.T) {
	t.Parallel()

	c := new(Container)
	Bind(c, Static(42))
	r := c.Freeze()

	out, err := TryResolve[int](r)
	assert.Equal(t, 42, out)
	assert.NoError(t, err)
}

func TestTryResolveContext(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "foo", "bar")

	c := new(Container)
	var outCtx interface{}
	Bind(c, Infallible(func(r Resolver) int {
		outCtx = r.Context().Value("foo")
		return 42
	}))
	r := c.Freeze()

	out, err := TryResolveContext[int](ctx, r)
	assert.Equal(t, 42, out)
	assert.NoError(t, err)
	assert.Equal(t, "bar", outCtx)
}

func TestTryResolveNamed(t *testing.T) {
	t.Parallel()

	c := new(Container)
	BindNamed(c, "foo", Static(42))
	r := c.Freeze()

	out, err := TryResolveNamed[int](r, "foo")
	assert.Equal(t, 42, out)
	assert.NoError(t, err)
}

func TestResolve(t *testing.T) {
	t.Parallel()

	c := new(Container)
	Bind(c, Static(42))
	r := c.Freeze()

	out := Resolve[int](r)
	assert.Equal(t, 42, out)
	assert.Panics(t, func() { Resolve[string](r) })
}

func TestResolveContext(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "foo", "bar")

	c := new(Container)
	var outCtx interface{}
	Bind(c, Infallible(func(r Resolver) int {
		outCtx = r.Context().Value("foo")
		return 42
	}))
	r := c.Freeze()

	out := ResolveContext[int](ctx, r)
	assert.Equal(t, 42, out)
	assert.Equal(t, "bar", outCtx)

	assert.Panics(t, func() { ResolveContext[string](ctx, r) })
}

func TestResolveNamed(t *testing.T) {
	t.Parallel()

	c := new(Container)
	BindNamed(c, "foo", Static(42))
	r := c.Freeze()

	out := ResolveNamed[int](r, "foo")
	assert.Equal(t, 42, out)
	assert.Panics(t, func() { ResolveNamed[string](r, "bar") })
}

func TestResolveNamedContext(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "foo", "bar")

	c := new(Container)
	var outFoo interface{}
	BindNamed(c, "fizz", func(r Resolver) (int, error) {
		outFoo = r.Context().Value("foo")
		return 123, nil
	})
	r := c.Freeze()

	out := ResolveNamedContext[int](ctx, r, "fizz")
	assert.Equal(t, 123, out)
	assert.Equal(t, "bar", outFoo)
	assert.Panics(t, func() { ResolveNamedContext[string](ctx, r, "bar") })
}
