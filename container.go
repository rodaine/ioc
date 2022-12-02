package ioc

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

const anonymous = ""

// Container is an inversion-of-control container. Providers for type
// constructors can be bound to the container using Bind and BindNamed. Once
// all providers have been attached, Freeze can be called to obtain a Resolver.
// Container is NOT thread-safe, but the resulting Resolver and associated
// functions are safe.
type Container struct {
	_ noCopy

	lookup map[typeName]providerFunc
	frozen bool
}

// Freeze prevents any more providers from being bound to the Container and
// returns a Resolver that can be used with the TryResolve* and Resolve*
// functions to produce values from the container. Calls to Bind or BindNamed
// after Freeze has been called will result in a panic. Freeze is idempotent
// and can be called multiple times safely.
func (c *Container) Freeze() Resolver {
	c.frozen = true
	return Resolver{c: c}
}

// BindNamed associates a ProviderFunc with the specified name and type. Note
// that type aliases (type Foo = Bar) are treated as the same type.
func BindNamed[T any](c *Container, name string, fn ProviderFunc[T]) {
	if c.lookup == nil {
		c.lookup = map[typeName]providerFunc{}
	}

	if c.frozen {
		panic("ioc.Resolver is frozen; no new providers may be bound")
	}

	c.lookup[newTypeName[T](name)] = fn.provide
}

// Bind associates a ProviderFunc with the specified type anonymously. It is
// equivalent to calling BindNamed with an empty name argument.
func Bind[T any](c *Container, fn ProviderFunc[T]) {
	BindNamed[T](c, anonymous, fn)
}
