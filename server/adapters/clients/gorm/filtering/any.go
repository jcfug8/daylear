package filtering

// Any represents a slice of any type.
type Any[T any] []T

func (Any[T]) isAny() {}

func appendAny[T any](items any, item T) (Any[T], bool) {
	switch items := items.(type) {
	case Any[T]:
		return append(items, item), true
	case nil:
		return Any[T]{item}, true
	default:
		return nil, false
	}
}
