package ioc

import "context"

// Resolver can be used once a Container is frozen (via Container.Freeze) to
// initialize values from the bound providers via the TryResolve* and Resolve*
// functions. Resolver and the associated functions are thread-safe.
type Resolver struct {
	c         *Container
	ctx       context.Context
	resolving []typeName
}

// Context returns the context.Context threaded through the resolve graph. By
// default, this value is equivalent to context.Background, but can be
// overridden by calling executing TryResolveContext* and similar functions.
func (r Resolver) Context() context.Context {
	if r.ctx != nil {
		return r.ctx
	}

	return context.Background()
}

func (r Resolver) startResolving(ctx context.Context, tn typeName) (Resolver, error) {
	chain := make([]typeName, len(r.resolving), 1+len(r.resolving))
	copy(chain, r.resolving)
	chain = append(chain, tn)

	for _, n := range r.resolving {
		if n == tn {
			return r, CircularDependencyError(chain)
		}
	}

	return Resolver{
		c:         r.c,
		ctx:       ctx,
		resolving: chain,
	}, nil
}

// TryResolveNamedContext will attempt to resolve a value for type T with the
// specified name. The provided context.Context will be passed to the target
// ProviderFunc via the Resolver.Context method.
//
// An error is returned if a provider cannot be found for the specified type
// and name (MissingProviderError), if there is dependency cycle in resolving
// (CircularDependencyError), or if the provider returns an error.
func TryResolveNamedContext[T any](ctx context.Context, r Resolver, name string) (value T, err error) {
	tn := newTypeName[T](name)

	provider := r.c.lookup[tn]
	if provider == nil {
		return value, MissingProviderError(tn)
	}

	r, err = r.startResolving(ctx, tn)
	if err != nil {
		return value, err
	}

	v, err := provider(r)
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
func TryResolveContext[T any](ctx context.Context, r Resolver) (value T, err error) {
	return TryResolveNamedContext[T](ctx, r, anonymous)
}

// TryResolveNamed will attempt to resolve a value for type T with the
// specified name.
//
// An error is returned if a provider cannot be found for the specified type
// and name (MissingProviderError), if there is dependency cycle in resolving
// (CircularDependencyError), or if the provider returns an error.
func TryResolveNamed[T any](r Resolver, name string) (T, error) {
	return TryResolveNamedContext[T](r.ctx, r, name)
}

// TryResolve will attempt to resolve a value for type T.
//
// An error is returned if a provider cannot be found for the specified type
// and name (MissingProviderError), if there is dependency cycle in resolving
// (CircularDependencyError), or if the provider returns an error.
func TryResolve[T any](r Resolver) (T, error) {
	return TryResolveContext[T](r.ctx, r)
}

// ResolveNamedContext behaves like TryResolveNamedContext, but panics in the
// event of an error resolving a value.
func ResolveNamedContext[T any](ctx context.Context, r Resolver, name string) T {
	value, err := TryResolveNamedContext[T](ctx, r, name)
	if err != nil {
		panic(err)
	}

	return value
}

// ResolveNamed behaves like TryResolveNamed, but panics in the event of an
// error resolving a value.
func ResolveNamed[T any](r Resolver, name string) T {
	return ResolveNamedContext[T](r.ctx, r, name)
}

// ResolveContext behaves like TryResolveContext, but panics in the event of an
// error resolving a value.
func ResolveContext[T any](ctx context.Context, r Resolver) T {
	return ResolveNamedContext[T](ctx, r, anonymous)
}

// Resolve behaves like TryResolve, but panics in the event of an error
// resolving a value.
func Resolve[T any](r Resolver) T {
	return ResolveContext[T](r.ctx, r)
}
