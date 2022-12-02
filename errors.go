package ioc

import (
	"fmt"
	"strings"
)

// CircularDependencyError is returned when calling a TryResolve* function
// results in a cycle of ProviderFunc.
type CircularDependencyError []typeName

func (err CircularDependencyError) Error() string {
	switch len(err) {
	case 0:
		return "circular dependency encountered"
	case 1:
		return fmt.Sprintf("circular dependency encountered resolving %v", err[0])
	default:
		sb := &strings.Builder{}
		_, _ = fmt.Fprintf(sb, "circular dependency encountered resolving %v:", err[0])

		for _, tn := range err[1:] {
			_, _ = fmt.Fprintf(sb, "\n- depends on %v", tn)
		}

		return sb.String()
	}
}

// MissingProviderError is returned when calling a TryResolve* function cannot
// find an associated ProviderFunc with the given type or name.
type MissingProviderError typeName

func (err MissingProviderError) Error() string {
	return fmt.Sprintf("missing provider for %v", typeName(err))
}

var (
	_ error = CircularDependencyError{}
	_ error = MissingProviderError{}
)
