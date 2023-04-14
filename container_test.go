package ioc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type someStruct struct {
	s string
}

func testBind[T any](t *testing.T, x T) {
	t.Helper()
	t.Run(fmt.Sprintf("%T", x), func(t *testing.T) {
		t.Parallel()
		c := new(Container)
		Bind(c, Static(x))
		assert.Equal(t, x, Resolve[T](c))
	})
}

func TestBind(t *testing.T) {
	t.Parallel()
	testBind(t, true)
	testBind(t, "foobar")
	testBind(t, 8)
	testBind(t, int8(16))
	testBind(t, int16(32))
	testBind(t, int32(64))
	testBind(t, int64(128))
	testBind(t, uint(8))
	testBind(t, uint8(16))
	testBind(t, uint16(32))
	testBind(t, uint32(64))
	testBind(t, uint64(128))
	testBind(t, uintptr(0xDEAD_BEEF))
	testBind(t, byte(16))
	testBind(t, 'x')
	testBind(t, float32(1.23))
	testBind(t, 4.56)
	testBind(t, complex64(1+2i))
	testBind(t, 3+4i)

	testBind(t, [3]int{1, 2, 3})
	testBind(t, []int{4, 5, 6})
	testBind(t, [2][3]int{{1, 2, 3}, {4, 5, 6}})
	testBind(t, map[string]int{"foo": 1, "bar": 2})

	type funcLocal struct{ x int }
	testBind(t, funcLocal{x: 1})
	testBind(t, someStruct{s: "foo"})
	testBind(t, struct{}{})
	testBind(t, struct{ y int }{y: 1})

	testBind(t, make(chan int))
	testBind(t, make(chan struct{}, 3))
	testBind(t, make(<-chan bool))
	testBind(t, make(chan<- string, 4))

	n := 123
	m := &n
	testBind(t, m)
	testBind(t, &m)
	testBind(t, &[]string{"alpha", "beta"})
	testBind(t, &someStruct{s: "bar"})
}

func TestBind_Alias(t *testing.T) {
	t.Parallel()

	c := new(Container)
	Bind(c, Static[uint8](2))
	assert.Equal(t, byte(2), Resolve[byte](c),
		"should resolve to the same provider")
}

func TestBind_Funcs(t *testing.T) {
	t.Parallel()

	fn := func(x int) int { return x }

	c := new(Container)
	Bind(c, Static(fn))

	assert.Equal(t,
		reflect.ValueOf(fn).Pointer(),
		reflect.ValueOf(Resolve[func(int) int](c)).Pointer())

	c = new(Container)
	buf := &bytes.Buffer{}
	Bind(c, Static(buf.WriteString))

	assert.Equal(t,
		reflect.ValueOf(buf.WriteString).Pointer(),
		reflect.ValueOf(Resolve[func(string) (int, error)](c)).Pointer())
}

func TestBind_Interfaces(t *testing.T) {
	t.Parallel()

	c := new(Container)
	Bind(c, Static[any]("foobar"))
	assert.Equal(t, "foobar", Resolve[any](c))

	c = new(Container)
	buf := &bytes.Buffer{}
	Bind(c, Static[io.Writer](buf))
	assert.Equal(t, buf, Resolve[io.Writer](c))
}

func TestBindNamed(t *testing.T) {
	t.Parallel()

	x := 123
	y := 456

	c := new(Container)
	BindNamed[int](c, "x", Static(x))
	BindNamed[int](c, "y", Static(y))

	assert.Equal(t, x, ResolveNamed[int](c, "x"))
	assert.Equal(t, y, ResolveNamed[int](c, "y"))
}

func TestContainer_Freeze(t *testing.T) {
	t.Parallel()

	c := new(Container)
	assert.NotPanics(t, func() { Bind(c, Static(123)) })

	c.Freeze()
	assert.Panics(t, func() { Bind(c, Static("foo")) })
}

func TestContainer_Context(t *testing.T) {
	t.Parallel()

	c := new(Container)
	assert.Equal(t, context.Background(), c.Context())
}

func TestContainer_Extend(t *testing.T) {
	t.Parallel()

	c := new(Container)
	Bind(c, Static(123))
	c.Freeze()

	ext := c.Extend()
	Bind(ext, Static("foo"))

	assert.Equal(t, 123, Resolve[int](ext))
	assert.Equal(t, "foo", Resolve[string](ext))
}
