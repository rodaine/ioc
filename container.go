package ioc

import (
	"context"
	"sync/atomic"
)

type noCopy struct{}

func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

const anonymous = ""

// Container is an inversion-of-control container. Providers for type
// constructors can be bound to the container using Bind and BindNamed. Once
// all providers have been attached, Freeze can be called to obtain a Resolver.
// The Container and associated functions are thread-safe.
type Container struct {
	_         noCopy
	parent    *Container
	providers syncMap[typeName, providerFunc]
	frozen    atomic.Bool
	resolving typeName
	ctx       context.Context
}

// Freeze prevents any more providers from being bound to the Container and
// returns a Resolver that can be used with the TryResolve* and Resolve*
// functions to produce values from the container. Calls to Bind or BindNamed
// after Freeze has been called will result in a panic. Freeze is idempotent
// and can be called multiple times safely.
func (c *Container) Freeze() {
	c.frozen.Store(true)
}

// Extend creates a new Container with the current container as its parent. This
// allows for creating non-destructive overrides to the bindings on the parent
// Container. A frozen Container can be extended to permit further bindings. A
// Container can be extended multiple times simultaneously.
//
// This method is primarily useful for creating a base, shared Container that is
// then extended with more domain-specific bindings to scope access.
func (c *Container) Extend() *Container {
	return &Container{
		parent: c,
		ctx:    c.ctx,
	}
}

// Context returns the context.Context threaded through the resolve graph. By
// default, this value is equivalent to context.Background, but can be
// overridden by executing TryResolveContext* and similar functions.
func (c *Container) Context() context.Context {
	if c.ctx != nil {
		return c.ctx
	}
	return context.Background()
}

func (c *Container) startResolving(ctx context.Context, name typeName) (*Container, error) {
	for rc := c; rc != nil; rc = rc.parent {
		if rc.resolving == name {
			return nil, CircularDependencyError(c.resolvingChain(name))
		}
	}

	resolver := &Container{
		parent:    c,
		ctx:       ctx,
		resolving: name,
	}
	resolver.Freeze()

	return resolver, nil
}

func (c *Container) resolvingChain(tn typeName) []typeName {
	chain := []typeName{tn}
	for rc := c; rc != nil; rc = rc.parent {
		if rc.resolving != (typeName{}) {
			chain = append(chain, rc.resolving)
		}
	}
	for i, j := 0, len(chain)-1; i < j; i, j = i+1, j-1 {
		chain[i], chain[j] = chain[j], chain[i]
	}
	return chain
}

func (c *Container) lookup(name typeName) (providerFunc, error) {
	for rc := c; rc != nil; rc = rc.parent {
		if provider, ok := rc.providers.Load(name); ok {
			return provider, nil
		}
	}
	return nil, MissingProviderError(name)
}

// BindNamed associates a ProviderFunc with the specified name and type. Note
// that type aliases (type Foo = Bar) are treated as the same type.
func BindNamed[T any](c *Container, name string, fn ProviderFunc[T]) {
	if c.frozen.Load() {
		panic("ioc.Container is frozen; no new providers may be bound")
	}

	c.providers.Store(newTypeName[T](name), fn.provide)
}

// Bind associates a ProviderFunc with the specified type anonymously. It is
// equivalent to calling BindNamed with an empty name argument.
func Bind[T any](c *Container, fn ProviderFunc[T]) {
	BindNamed[T](c, anonymous, fn)
}
