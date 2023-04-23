package timtab

func NewSet[T comparable](ts ...T) *Set[T] {
	s := &Set[T]{
		values: make(map[T]struct{}),
	}
	for _, t := range ts {
		s.Insert(t)
	}
	return s
}

type Set[T comparable] struct {
	values map[T]struct{}
}

func (s *Set[T]) Insert(t T) {
	s.values[t] = struct{}{}
}

func (s *Set[T]) Contains(t T) bool {
	_, ok := s.values[t]
	return ok
}

func (s *Set[T]) Intersects(os *Set[T]) bool {
	for t := range s.values {
		if os.Contains(t) {
			return true
		}
	}
	return false
}
