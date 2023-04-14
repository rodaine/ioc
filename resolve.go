package ioc

import "context"

// TryResolveNamedContext will attempt to resolve a value for type T with the
// specified name. The provided context.Context will be passed to the target
// ProviderFunc via the Resolver.Context method.
//
// An error is returned if a provider cannot be found for the specified type
// and name (MissingProviderError), if there is dependency cycle in resolving
// (CircularDependencyError), or if the provider returns an error.
func TryResolveNamedContext[T any](ctx context.Context, container *Container, name string) (value T, err error) {
	tname := newTypeName[T](name)

	provider, err := container.lookup(tname)
	if err != nil {
		return value, err
	}

	container, err = container.startResolving(ctx, tname)
	if err != nil {
		return value, err
	}

	v, err := provider(container)
	if err != nil {
		return value, err
	}

	return v.(T), nil
}

// TryResolveContext will attempt to resolve a value for type T. The provided
// context.Context will be passed to the target ProviderFunc via the
// Resolver.Context method.
//
// An error is returned if a provider cannot be found for the specified type
// and name (MissingProviderError), if there is dependency cycle in resolving
// (CircularDependencyError), or if the provider returns an error.
func TryResolveContext[T any](ctx context.Context, c *Container) (value T, err error) {
	return TryResolveNamedContext[T](ctx, c, anonymous)
}

// TryResolveNamed will attempt to resolve a value for type T with the
// specified name.
//
// An error is returned if a provider cannot be found for the specified type
// and name (MissingProviderError), if there is dependency cycle in resolving
// (CircularDependencyError), or if the provider returns an error.
func TryResolveNamed[T any](c *Container, name string) (T, error) {
	return TryResolveNamedContext[T](c.ctx, c, name)
}

// TryResolve will attempt to resolve a value for type T.
//
// An error is returned if a provider cannot be found for the specified type
// and name (MissingProviderError), if there is dependency cycle in resolving
// (CircularDependencyError), or if the provider returns an error.
func TryResolve[T any](c *Container) (T, error) {
	return TryResolveContext[T](c.ctx, c)
}

// ResolveNamedContext behaves like TryResolveNamedContext, but panics in the
// event of an error resolving a value.
func ResolveNamedContext[T any](ctx context.Context, c *Container, name string) T {
	value, err := TryResolveNamedContext[T](ctx, c, name)
	if err != nil {
		panic(err)
	}

	return value
}

// ResolveNamed behaves like TryResolveNamed, but panics in the event of an
// error resolving a value.
func ResolveNamed[T any](c *Container, name string) T {
	return ResolveNamedContext[T](c.ctx, c, name)
}

// ResolveContext behaves like TryResolveContext, but panics in the event of an
// error resolving a value.
func ResolveContext[T any](ctx context.Context, c *Container) T {
	return ResolveNamedContext[T](ctx, c, anonymous)
}

// Resolve behaves like TryResolve, but panics in the event of an error
// resolving a value.
func Resolve[T any](c *Container) T {
	return ResolveContext[T](c.ctx, c)
}
