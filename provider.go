package ioc

import "sync"

type providerFunc func(*Container) (any, error)

// ProviderFunc is a function that creates a new instance of T, returning an
// error if the operation cannot be performed. A ProviderFunc is bound to a
// Container via Bind or BindNamed.
//
// A ProviderFunc may resolve other bound providers in the process of
// initializing the desired return value. It is a best practice to use the
// TryResolve* functions (instead of Resolve*) within ProviderFunc
// implementations to allow for consuming code to control the failure behavior.
type ProviderFunc[T any] func(c *Container) (T, error)

func (fn ProviderFunc[T]) provide(c *Container) (any, error) { return fn(c) }

// Singleton wraps fn, returning a new ProviderFunc with the same signature,
// but the returned values always remain the same. Singleton is useful for
// defining a provider to a value that should be shared, such as
// network/database clients, loggers, or other thread-safe utilities.
func Singleton[T any](provider ProviderFunc[T]) ProviderFunc[T] {
	var value T
	var err error
	once := &sync.Once{}

	return func(c *Container) (T, error) {
		once.Do(func() { value, err = provider(c) })
		return value, err
	}
}

// Static creates a ProviderFunc that always returns (v, nil). Static is useful
// where the value does not have other dependencies that need to be resolved
// before consumption.
func Static[T any](v T) ProviderFunc[T] {
	return func(*Container) (T, error) { return v, nil }
}

// Infallible is a helper that converts a func(c Resolver) T to a ProviderFunc
// where an error is never returned.
func Infallible[T any](fn func(c *Container) T) ProviderFunc[T] {
	return func(c *Container) (T, error) { return fn(c), nil }
}

var _ providerFunc = (ProviderFunc[any])(nil).provide
