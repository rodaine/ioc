package ioc

import "sync"

type providerFunc func(Resolver) (any, error)

// ProviderFunc is a function that creates a new instance of T, returning an
// error if the operation cannot be performed. A ProviderFunc is bound to a
// Resolver via Bind or BindNamed.
//
// A ProviderFunc may resolve other bound providers in the process of
// initializing the desired return value. It is a best practice to use the
// TryResolve* functions (instead of Resolve*) within ProviderFunc
// implementations to allow for consuming code to control the failure behavior.
type ProviderFunc[T any] func(r Resolver) (T, error)

func (fn ProviderFunc[T]) provide(r Resolver) (any, error) { return fn(r) }

// Singleton wraps fn, returning a new ProviderFunc with the same signature,
// but the returned values always remain the same. Singleton is useful for
// defining a provider to a value that should be shared, such as
// network/database clients, loggers, or other thread-safe utilities.
func Singleton[T any](fn ProviderFunc[T]) ProviderFunc[T] {
	var v T
	var err error
	once := &sync.Once{}

	return func(r Resolver) (T, error) {
		once.Do(func() { v, err = fn(r) })
		return v, err
	}
}

// Static creates a ProviderFunc that always returns (v, nil). Static is useful
// where the value does not have other dependencies that need to be resolved
// before consumption.
func Static[T any](v T) ProviderFunc[T] {
	return func(Resolver) (T, error) { return v, nil }
}

// Infallible is a helper that converts a func(c Resolver) T to a ProviderFunc
// where an error is never returned.
func Infallible[T any](fn func(c Resolver) T) ProviderFunc[T] {
	return func(r Resolver) (T, error) { return fn(r), nil }
}

var _ providerFunc = (ProviderFunc[any])(nil).provide
