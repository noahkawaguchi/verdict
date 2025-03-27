package utils

type set[T comparable] struct{ m map[T]struct{} }

// NewSet creates an unordered collection of unique, comparable elements.
func NewSet[T comparable](items ...T) *set[T] {
	s := &set[T]{make(map[T]struct{})}
	for _, item := range items {
		s.m[item] = struct{}{}
	}
	return s
}

// Len returns the number of elements in the set (already unique).
func (s *set[T]) Len() int { return len(s.m) }
